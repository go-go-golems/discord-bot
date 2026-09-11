// Source-only launcher; no package installation or environment credentials.
// Usage: bun probe.ts MOCK_CHECKOUT RUN_DIRECTORY
import { mkdir, writeFile, rename } from "node:fs/promises";
import { resolve } from "node:path";
const [checkout, directory] = Bun.argv.slice(2);
if (!checkout || !directory) throw new Error("provide mock checkout and run directory");
const { SlackMock } = await import(resolve(checkout, "src/server.ts"));
await mkdir(directory, { recursive: true, mode: 0o700 });
const mock = await SlackMock.start({
  host: "127.0.0.1", port: 0, subscribedEvents: ["app_mention"],
  manifest: { features: { slash_commands: [{ command: "/golem-ping", description: "Probe" }] } },
});
const deadline = setTimeout(() => { console.error("probe deadline exceeded"); process.exit(1); }, 45000);
const write = async (name: string, value: unknown) => {
  await writeFile(resolve(directory, name + ".tmp"), JSON.stringify(value, null, 2), { mode: 0o600 });
  await rename(resolve(directory, name + ".tmp"), resolve(directory, name));
};
try {
  await write("config.json", { apiURL: mock.apiUrl, botToken: mock.bot.token, appToken: mock.bot.appToken,
    teamID: mock.team.id, userID: mock.bot.userId });
  console.log(JSON.stringify({ phase: "ready", port: mock.port }));
  await mock.waitForConnection(20000);
  const mention = await mock.postMessage({ channel: "general", user: "alice", text: `<@${mock.bot.userId}> probe-mention` });
  const reply = await mock.waitForMessage({ text: "probe-thread-reply", thread_ts: mention.ts }, { timeoutMs: 10000 });
  const command = await mock.slashCommand({ command: "/golem-ping", user: "alice", channel: "general", text: "probe-command" });
  const ephemeral = await mock.waitForMessage({ text: "probe-private-reply", ephemeral: true }, { timeoutMs: 10000 });
  const deliveries = mock.deliveries().map((d: any) => ({ name: d.name, acked: d.acked, attempts: d.attempts }));
  const result = { threadPreserved: reply.thread_ts === mention.ts,
    ephemeralRecipient: ephemeral.ephemeral_user === "U0ALICE000",
    publicReplyCount: mock.findMessages({ text: "probe-private-reply" }).length,
    emptyCommandAck: command.ack.payload == null, deliveries };
  if (!result.threadPreserved || !result.ephemeralRecipient || result.publicReplyCount !== 0 || !result.emptyCommandAck || deliveries.length !== 2 || deliveries.some((d: any) => !d.acked)) throw new Error("mock assertions failed");
  await write("result.json", result);
  while (mock.connectionCount !== 0) await Bun.sleep(25);
  await write("shutdown.json", { connections: mock.connectionCount });
  console.log(JSON.stringify({ phase: "passed", ...result }));
} finally {
  clearTimeout(deadline);
  await mock.stop();
}
