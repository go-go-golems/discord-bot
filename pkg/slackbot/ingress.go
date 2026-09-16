package slackbot

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Dispatcher interface {
	Dispatch(context.Context, Invocation, Responder) error
}
type Acknowledger interface {
	Ack(context.Context, string, bool) error
}

// Envelope is a detached transport snapshot. EventID identifies a Slack event;
// ID identifies the socket delivery. Bot excludes bot-authored messages.
type Envelope struct {
	ID, EventID, AppID string
	Bot                bool
	Invocation         Invocation
	Responder          Responder
}
type Admission string

const (
	Accepted  Admission = "accepted"
	Duplicate Admission = "duplicate"
	Dropped   Admission = "dropped"
	Busy      Admission = "busy"
	Closed    Admission = "closed"
)

type IngressOptions struct {
	TeamID, AppID, SelfUserID string
	AllowedChannels           []string
	Capacity, DedupeCapacity  int
	DedupeTTL                 time.Duration
}
type IngressStats struct{ Accepted, Duplicate, Dropped, Busy uint64 }

// Ingress owns bounded in-memory admission. It acknowledges without entering JS.
// ACK does not guarantee durable execution. Dispatchers must honor cancellation.
type Ingress struct {
	ctx      context.Context
	cancel   context.CancelFunc
	group    errgroup.Group
	queue    chan Envelope
	done     chan struct{}
	failures chan error
	opts     IngressOptions
	mu       sync.Mutex
	seen     map[string]time.Time
	stats    IngressStats
}

func NewIngress(ctx context.Context, opts IngressOptions, dispatcher Dispatcher) (*Ingress, error) {
	if opts.TeamID == "" || opts.AppID == "" || dispatcher == nil {
		return nil, errors.New("ingress requires team, app and dispatcher")
	}
	if opts.Capacity <= 0 || opts.DedupeCapacity <= 0 || opts.DedupeTTL <= 0 {
		return nil, errors.New("ingress capacities and TTL must be positive")
	}
	opts.AllowedChannels = append([]string(nil), opts.AllowedChannels...)
	ctx, cancel := context.WithCancel(ctx)
	p := &Ingress{ctx: ctx, cancel: cancel, queue: make(chan Envelope, opts.Capacity), done: make(chan struct{}), failures: make(chan error, opts.Capacity), opts: opts, seen: map[string]time.Time{}}
	p.group.Go(func() error {
		defer close(p.done)
		defer close(p.failures)
		for {
			select {
			case <-ctx.Done():
				return nil
			case e := <-p.queue:
				if ctx.Err() != nil {
					return nil
				}
				if err := dispatcher.Dispatch(ctx, e.Invocation, e.Responder); err != nil && ctx.Err() == nil {
					select {
					case p.failures <- err:
					default:
					} // bounded diagnostics, never stall receipt
				}
			}
		}
	})
	return p, nil
}
func (p *Ingress) Close(ctx context.Context) error {
	p.cancel()
	select {
	case <-p.done:
		return p.group.Wait()
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (p *Ingress) Errors() <-chan error { return p.failures }
func (p *Ingress) Stats() IngressStats  { p.mu.Lock(); defer p.mu.Unlock(); return p.stats }
func (p *Ingress) Admit(ctx context.Context, e Envelope, ack Acknowledger) (Admission, error) {
	if e.ID == "" || ack == nil {
		return Dropped, errors.New("envelope ID and acknowledger are required")
	}
	decision := p.admit(e, time.Now())
	// Transport owns the ACK context and its deadline. The worker uses p.ctx instead.
	// View submissions choose their response-bearing ACK from JavaScript. The
	// transport receipt is carried on the invocation and enforces single-use and
	// deadline rules; admitting it here must not accept the submission early.
	if e.Invocation.Interaction == nil {
		if err := ack.Ack(ctx, e.ID, decision == Busy); err != nil {
			return decision, errors.Wrap(err, "acknowledge envelope")
		}
	}
	if e.Invocation.Interaction != nil && e.Invocation.Interaction.Ack == nil {
		return decision, errors.New("interactive invocation requires an acknowledger")
	}
	if e.Invocation.Interaction == nil {
		return decision, nil
	}
	if decision == Busy || decision == Dropped || decision == Duplicate || decision == Closed {
		// Leave explicit interactive envelopes unacknowledged when they cannot be
		// admitted; Slack may retry them and the handler never accepts data it did
		// not process.
		return decision, nil
	}
	return decision, nil
}
func (p *Ingress) admit(e Envelope, now time.Time) Admission {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ctx.Err() != nil {
		return Closed
	}
	allowed := len(p.opts.AllowedChannels) == 0
	for _, id := range p.opts.AllowedChannels {
		if id == e.Invocation.ChannelID {
			allowed = true
			break
		}
	}
	if !allowed || e.Bot || e.Invocation.UserID == p.opts.SelfUserID || e.AppID != p.opts.AppID || e.Invocation.TeamID != p.opts.TeamID || e.Invocation.Validate() != nil {
		p.stats.Dropped++
		return Dropped
	}
	key := "envelope:" + e.ID
	if e.Invocation.Event != "" {
		if e.EventID == "" {
			p.stats.Dropped++
			return Dropped
		}
		key = "event:" + e.Invocation.TeamID + ":" + e.EventID
	}
	for k, expiry := range p.seen {
		if !expiry.After(now) {
			delete(p.seen, k)
		}
	}
	if _, exists := p.seen[key]; exists {
		p.stats.Duplicate++
		return Duplicate
	}
	// Do not evict a live key to make room: preserve dedupe guarantees until TTL.
	if len(p.seen) >= p.opts.DedupeCapacity {
		p.stats.Busy++
		return Busy
	}
	select {
	case p.queue <- e:
		p.seen[key] = now.Add(p.opts.DedupeTTL)
		p.stats.Accepted++
		return Accepted
	default:
		p.stats.Busy++
		return Busy
	}
}
