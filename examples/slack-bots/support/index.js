const { defineBot } = require("slack"),
  ui = require("slack/ui");
module.exports = defineBot(({ configure, command, event }) => {
  configure({
    name: "support",
    description: "Support drafts and Slack conversation workflows",
    scopes: ["channels:history", "channels:read", "channels:join"],
  });
  command(
    "/support-ticket",
    { description: "Draft a support ticket" },
    async (ctx) => {
      const topic = ctx.text.trim();
      if (!topic) return { text: "Usage: /support-ticket topic" };
      await ctx.reply(
        ui
          .message("Drafted a support ticket for " + topic)
          .block(ui.header("Support Ticket Draft"))
          .block(ui.section(ui.plain(topic)))
          .block(
            ui.actions(
              "handbook",
              ui.linkButton(
                "handbook",
                "Support Handbook",
                "https://example.com/support-handbook",
              ),
            ),
          )
          .build(),
      );
      await ctx.slack.messages.ephemeral({
        channel: ctx.channelId,
        user: ctx.userId,
        text: "A follow-up reminder was also sent.",
      });
    },
  );
  command("/support-status", { description: "Support status" }, async () => ({
    text: "Support systems nominal.",
  }));
  command(
    "/support-fetch-thread",
    { description: "Fetch replies: parent timestamp" },
    async (ctx) => {
      const ts = ctx.text.trim();
      if (!ts) return { text: "Usage: /support-fetch-thread parent-ts" };
      const result = await ctx.slack.conversations.replies({
        channel: ctx.channelId,
        ts,
        limit: 100,
      });
      return { text: JSON.stringify(result.messages || []).slice(0, 3900) };
    },
  );
  command(
    "/support-start-thread",
    { description: "Start a named conversation thread" },
    async (ctx) => {
      if (!ctx.text.trim())
        return { text: "Usage: /support-start-thread topic" };
      const ref = await ctx.slack.messages.post({
        channelId: ctx.channelId,
        text: ctx.text,
      });
      await ctx.slack.messages.post({
        channelId: ctx.channelId,
        threadTs: ref.ts,
        text: "Support discussion starts here.",
      });
      return { text: "Started support thread " + ref.ts };
    },
  );
  // Slack threads have no independent join/leave membership. Channel membership is explicit.
  command(
    "/support-join-channel",
    { description: "Join a public Slack channel by ID" },
    async (ctx) => {
      const result = await ctx.slack.conversations.join({
        channel: ctx.text.trim() || ctx.channelId,
      });
      return {
        text:
          "Joined channel " + (result.channel ? result.channel.id : ctx.text),
      };
    },
  );
  command(
    "/support-leave-channel",
    { description: "Leave the current Slack channel" },
    async (ctx) => {
      await ctx.reply({ text: "Leaving channel." });
      await ctx.slack.conversations.leave({ channel: ctx.channelId });
    },
  );
  event("message", async (ctx) => {
    if (ctx.text.trim() === "!support")
      await ctx.reply({ text: "Support bot received your message trigger." });
  });
});
