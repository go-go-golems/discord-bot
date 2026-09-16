const { defineBot } = require("slack"),
  ui = require("slack/ui");
async function archive(ctx, channel, thread, limit, before) {
  const info = await ctx.slack.conversations.info({ channel });
  const messages = [],
    seen = new Set(),
    cursors = new Set();
  let cursor = "";
  do {
    const params = { channel, limit: Math.min(100, limit - messages.length) };
    if (cursor) params.cursor = cursor;
    if (before) params.latest = before;
    if (thread) params.ts = thread;
    const page = await (thread
      ? ctx.slack.conversations.replies(params)
      : ctx.slack.conversations.history(params));
    for (const message of page.messages || []) {
      if (!seen.has(message.ts)) {
        messages.push(message);
        seen.add(message.ts);
        if (messages.length >= limit) break;
      }
    }
    const next = (page.response_metadata || {}).next_cursor || "";
    if (next && cursors.has(next))
      throw Error(
        "Slack returned a repeated pagination cursor; archive aborted.",
      );
    if (next) cursors.add(next);
    cursor = next;
  } while (cursor && messages.length < limit);
  if (!messages.length) return "No messages found.";
  messages.sort((a, b) => a.ts.localeCompare(b.ts));
  const name = (info.channel || {}).name || channel;
  const link = await ctx.slack.messages.permalink({
    conversation_id: channel,
    message_ts: messages[0].ts,
  });
  const permalink = String(link.permalink || "").split("?")[0];
  let markdown =
    "# " +
    name +
    "\n\n" +
    (thread ? "Thread: " + thread + "\n\n" : "") +
    "Messages: " +
    messages.length +
    " (requested limit " +
    limit +
    ")\n";
  markdown =
    "---\nsource: slack\nteam_id: " +
    JSON.stringify(ctx.teamId) +
    "\nchannel_id: " +
    JSON.stringify(channel) +
    "\narchived_at: " +
    JSON.stringify(new Date().toISOString()) +
    "\nmessage_count: " +
    messages.length +
    "\n---\n\n" +
    markdown;
  for (const message of messages) {
    markdown +=
      "\n## " +
      new Date(Number(message.ts) * 1000).toISOString() +
      " · " +
      (message.user || message.bot_id || "unknown") +
      "\n\n" +
      (message.text || "") +
      "\n";
    if (message.edited) markdown += "\n*(edited)*\n";
    if (!message.text && message.blocks)
      markdown +=
        "\n```json\n" + JSON.stringify(message.blocks, null, 2) + "\n```\n";
    for (const attachment of message.attachments || [])
      markdown +=
        "\n" +
        (attachment.title || "") +
        "\n" +
        (attachment.text || attachment.fallback || "") +
        "\n" +
        (attachment.title_link || "") +
        "\n";
    if (permalink)
      markdown +=
        "\n[Open in Slack](" +
        permalink.replace(/p[0-9]+$/, "p" + message.ts.replace(".", "")) +
        ")\n";
    for (const file of message.files || [])
      markdown +=
        "\nAttachment: [" +
        (file.title || file.name || file.id) +
        "](" +
        (file.permalink || file.url_private || "") +
        ")\n";
    markdown += "\nSlack timestamp: " + message.ts + "\n";
  }
  await ctx.slack.files.upload({
    channel_id: channel,
    filename: name.replace(/[^a-zA-Z0-9_-]/g, "-") + "--archive.md",
    content: markdown,
  });
  return (
    "Archived " +
    messages.length +
    " messages" +
    (cursor ? "; requested limit reached." : ".")
  );
}
module.exports = defineBot(({ configure, command, shortcut, view }) => {
  configure({
    name: "archive-helper",
    description: "Paginated channel/thread Markdown archives",
    scopes: [
      "channels:history",
      "channels:read",
      "groups:history",
      "groups:read",
      "files:write",
    ],
    run: { fields: { defaultLimit: { type: "number", default: 500 } } },
  });
  command(
    "/archive-channel",
    { description: "Archive current channel: [limit] [before-ts]" },
    async (ctx) => {
      const args = ctx.text.trim().split(/\s+/);
      const limit = Number(args[0] || ctx.config.defaultLimit);
      if (!Number.isInteger(limit) || limit < 1 || limit > 10000)
        return { text: "Limit must be 1–10000." };
      return {
        text: await archive(ctx, ctx.channelId, "", limit, args[1] || ""),
      };
    },
  );
  shortcut(
    {
      callbackId: "archive-thread",
      name: "Archive Thread",
      description: "Export this message and its replies",
      type: "message",
    },
    async (ctx) => {
      await ctx.openModal(
        ui
          .modal("archive.thread", "Archive thread")
          .metadata(
            JSON.stringify({
              channel: ctx.channelId,
              thread: ctx.shortcut.message.thread_ts || ctx.shortcut.message.ts,
            }),
          )
          .input(
            "limit",
            "Maximum messages",
            ui.textInput("value").initial(String(ctx.config.defaultLimit)),
          )
          .submit("Export")
          .build(),
      );
    },
  );
  view("archive.thread", async (ctx) => {
    const limit = Number(ctx.values.text("limit", "value"));
    if (!Number.isInteger(limit) || limit < 1 || limit > 10000) {
      await ctx.ack.errors({ limit: "Use an integer from 1–10000." });
      return;
    }
    const target = JSON.parse(ctx.view.privateMetadata);
    await ctx.ack.accept();
    try {
      const text = await archive(ctx, target.channel, target.thread, limit, "");
      await ctx.slack.messages.ephemeral({
        channel: target.channel,
        user: ctx.userId,
        text,
      });
    } catch (err) {
      await ctx.slack.messages.ephemeral({
        channel: target.channel,
        user: ctx.userId,
        text: "Archive failed: " + err.message,
      });
    }
  });
});
