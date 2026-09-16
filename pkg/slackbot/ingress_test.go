package slackbot

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type blockedDispatcher struct {
	started chan struct{}
	once    sync.Once
}

var _ Dispatcher = (*blockedDispatcher)(nil)

func (d *blockedDispatcher) Dispatch(ctx context.Context, _ Invocation, _ Responder) error {
	d.once.Do(func() { close(d.started) })
	<-ctx.Done()
	return ctx.Err()
}

type ackRecorder struct {
	mu   sync.Mutex
	ids  []string
	busy []bool
}

var _ Acknowledger = (*ackRecorder)(nil)

func (a *ackRecorder) Ack(_ context.Context, id string, busy bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ids = append(a.ids, id)
	a.busy = append(a.busy, busy)
	return nil
}
func envelope(id, eventID string) Envelope {
	return Envelope{ID: id, EventID: eventID, AppID: "A", Invocation: Invocation{TeamID: "T", ChannelID: "C", UserID: "U", Event: "app_mention", TS: "1741234567.000001"}}
}
func TestIngressAckIndependentDedupeAndOverload(t *testing.T) {
	d := &blockedDispatcher{started: make(chan struct{})}
	a := &ackRecorder{}
	p, err := NewIngress(context.Background(), IngressOptions{TeamID: "T", AppID: "A", Capacity: 1, DedupeCapacity: 10, DedupeTTL: time.Minute}, d)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, p.Close(context.Background())) })
	got, err := p.Admit(context.Background(), envelope("one", "event1"), a)
	require.NoError(t, err)
	require.Equal(t, Accepted, got)
	select {
	case <-d.started:
	case <-time.After(time.Second):
		t.Fatal("dispatcher did not start")
	}
	got, err = p.Admit(context.Background(), envelope("two", "event1"), a)
	require.NoError(t, err)
	require.Equal(t, Duplicate, got)
	got, err = p.Admit(context.Background(), envelope("three", "event2"), a)
	require.NoError(t, err)
	require.Equal(t, Accepted, got)
	got, err = p.Admit(context.Background(), envelope("four", "event3"), a)
	require.NoError(t, err)
	require.Equal(t, Busy, got)
	require.Equal(t, []string{"one", "two", "three", "four"}, a.ids)
	require.True(t, a.busy[3])
	require.Equal(t, IngressStats{Accepted: 2, Duplicate: 1, Busy: 1}, p.Stats())
}
func TestIngressFiltersAndDedupeBound(t *testing.T) {
	d := &blockedDispatcher{started: make(chan struct{})}
	a := &ackRecorder{}
	p, err := NewIngress(context.Background(), IngressOptions{TeamID: "T", AppID: "A", SelfUserID: "BOT", AllowedChannels: []string{"C"}, Capacity: 2, DedupeCapacity: 1, DedupeTTL: time.Minute}, d)
	require.NoError(t, err)
	defer func() { _ = p.Close(context.Background()) }()
	for _, modify := range []func(*Envelope){func(e *Envelope) { e.Bot = true }, func(e *Envelope) { e.AppID = "other" }, func(e *Envelope) { e.Invocation.TeamID = "other" }, func(e *Envelope) { e.Invocation.ChannelID = "other" }, func(e *Envelope) { e.Invocation.UserID = "BOT" }, func(e *Envelope) { e.EventID = "" }} {
		e := envelope("drop", "e")
		modify(&e)
		got, err := p.Admit(context.Background(), e, a)
		require.NoError(t, err)
		require.Equal(t, Dropped, got)
	}
	got, err := p.Admit(context.Background(), envelope("one", "e1"), a)
	require.NoError(t, err)
	require.Equal(t, Accepted, got)
	got, err = p.Admit(context.Background(), envelope("two", "e2"), a)
	require.NoError(t, err)
	require.Equal(t, Busy, got)
	require.NoError(t, p.Close(context.Background()))
	got, err = p.Admit(context.Background(), envelope("closed", "e3"), a)
	require.NoError(t, err)
	require.Equal(t, Closed, got)
}

func TestDedupeTTLAndNonAdmission(t *testing.T) {
	// Exercise admission with an explicit clock so TTL assertions need no sleeps.
	d := &blockedDispatcher{started: make(chan struct{})}
	p, err := NewIngress(context.Background(), IngressOptions{TeamID: "T", AppID: "A", Capacity: 3, DedupeCapacity: 1, DedupeTTL: time.Minute}, d)
	require.NoError(t, err)
	defer func() { _ = p.Close(context.Background()) }()
	now := time.Now()
	require.Equal(t, Accepted, p.admit(envelope("a", "e"), now))
	require.Equal(t, Duplicate, p.admit(envelope("b", "e"), now.Add(time.Second)))
	require.Equal(t, Accepted, p.admit(envelope("c", "e"), now.Add(2*time.Minute)))
}
