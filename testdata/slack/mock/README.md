# Local SDK interoperability gate

This probe runs the actual pinned Go Slack SDK against the separately checked-out
Bun mock. It verifies identity, Socket Mode hello and two deliveries, exact ACKs
as observed by the mock, a threaded Web API post, an ephemeral response URL post,
and socket cancellation. It does not yet execute the application JavaScript host
or prove production permission enforcement, retry policy, or precise ACK timing.

Versions are recorded in `revisions.json`; prepare those dependencies before
running offline. The tested mock checkout is
`/home/manuel/code/others/slack/desplega-slack-mock`. No installation happens in
the probe. The mock launcher imports `src/server.ts` directly.

Run from the repository root, using a fresh run directory each time:

```sh
tmux new-session -d -s slack-sdk-probe \
  '/home/manuel/.bun/bin/bun testdata/slack/mock/probe.ts /home/manuel/code/others/slack/desplega-slack-mock /tmp/slack-sdk-probe-NEW; read'
# Wait for config.json, without printing its synthetic credentials.
go test -count=1 -v ./internal/slackprobe -run TestSDKInteroperability \
  -args -slack-probe-directory /tmp/slack-sdk-probe-NEW
tmux capture-pane -p -t slack-sdk-probe
```

On this workstation, use the ticket's `scripts/04-go-offline.sh` instead of `go`
to select the cached matching toolchain and disable downloads. The launcher has a
45-second outer deadline; the Go test has a 30-second operation context and waits
for explicit result/shutdown receipts. The server normally stops itself after
client disconnection. Clean the owned port with `lsof-who -p PORT -k`, then remove
the owned tmux session with `tmux kill-session -t slack-sdk-probe`.

The synthetic `config.json` stays outside Git with mode 0600. Do not archive it.
`result.json` and `shutdown.json` are sanitized evidence suitable for the ticket.
Writes use atomic rename so readiness cannot expose a partially written file.

Ordinary `go test ./...` explicitly skips the external-server test unless the
run-directory argument is supplied; it still runs the destination rejection
test. Therefore a successful ordinary suite is **not** a successful mock gate.
An explicitly selected nonexistent directory fails, rather than silently skipping.
Future CI must run this gate separately and fail when its prepared runtime is absent.

The Go HTTP and WebSocket dialers reject all non-literal/non-loopback addresses
before DNS or dialing. HTTP proxies are unset in the constructed transport,
redirects are rejected, and the same HTTP client handles response URLs. Only
synthetic credentials are loaded from the explicit file; no environment values
are read by the probe.
