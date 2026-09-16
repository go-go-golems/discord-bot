const { defineBot } = require("slack"),
  ui = require("slack/ui");
const topics = [
  {
    title: "Architecture",
    summary: "Bot wiring, handlers, and runtime layers.",
  },
  {
    title: "Testing",
    summary: "Go tests, JS runtime coverage, and smoke checks.",
  },
  { title: "Runbooks", summary: "How to start, sync, and operate the bot." },
];
function results(query) {
  const q = query.toLowerCase();
  const matches = topics.filter((t) =>
    (t.title + " " + t.summary).toLowerCase().includes(q),
  );
  return matches.length
    ? matches.map((t) => t.title + ": " + t.summary).join("\n")
    : "No results for " + query;
}
module.exports = defineBot(
  ({ configure, command, event, action, view, options }) => {
    configure({
      name: "ping",
      description:
        "Commands, native controls, modals, search and generated files",
      scopes: ["files:write"],
      run: {
        fields: {
          greeting: { type: "string", default: "pong", help: "Ping response" },
        },
      },
    });
    command(
      "/golem-ping",
      { description: "Check the development bot" },
      async (ctx) =>
        ui
          .message(ctx.config.greeting)
          .block(ui.header("Pong"))
          .block(ui.section(ui.plain("JavaScript handled this slash command.")))
          .block(
            ui.actions(
              "ping.controls",
              ui.button("ping.panel", "Open panel"),
              ui.linkButton(
                "ping.repo",
                "Project repo",
                "https://github.com/manuel/wesen",
              ),
            ),
          )
          .block(
            ui.actions(
              "ping.topics",
              ui.staticSelect("ping.topic", {
                placeholder: "Choose a topic",
                options: topics.map((t) =>
                  ui.option(t.title, t.title.toLowerCase()),
                ),
              }),
            ),
          )
          .build(),
    );
    command(
      "/golem-echo",
      { description: "Echo text with a private follow-up" },
      async (ctx) => {
        await ctx.reply({ text: ctx.text || "Usage: /golem-echo text" });
        await ctx.slack.messages.ephemeral({
          channel: ctx.channelId,
          user: ctx.userId,
          text: "Follow-up from JavaScript",
        });
      },
    );
    command(
      "/golem-feedback",
      { description: "Open a feedback form" },
      async (ctx) => {
        await ctx.openModal(
          ui
            .modal("ping.feedback", "Feedback")
            .input("summary", "Summary", ui.textInput("value").length(5, 100))
            .input(
              "details",
              "Details",
              ui.textInput("value").multiline().length(0, 500),
              { optional: true },
            )
            .submit("Submit")
            .build(),
        );
      },
    );
    command(
      "/golem-search",
      { description: "Search topics or open dynamic search" },
      async (ctx) => {
        if (ctx.text.trim()) return { text: results(ctx.text.trim()) };
        await ctx.openModal(
          ui
            .modal("ping.search", "Search topics")
            .input(
              "query",
              "Topic",
              ui.externalSelect("ping.search.query", { min_query_length: 0 }),
            )
            .submit("Search")
            .build(),
        );
      },
    );
    options("ping.search.query", async (ctx) => {
      const q = (ctx.query || "").toLowerCase();
      await ctx.ack.options(
        topics
          .filter((t) => t.title.toLowerCase().includes(q))
          .map((t) => ui.option(t.title, t.title.toLowerCase()))
          .concat([
            ui.option(
              "Custom: " + (ctx.query || "query"),
              (ctx.query || "custom").slice(0, 100),
            ),
          ]),
      );
    });
    view("ping.search", async (ctx) => {
      const option = ctx.values.all.query["ping.search.query"].selected_option;
      await ctx.ack.update(
        ui
          .modal("ping.search.results", "Search results")
          .block(ui.section(ui.plain(results(option.value))))
          .build(),
      );
    });
    view("ping.feedback", async (ctx) => {
      await ctx.ack.update(
        ui
          .modal("ping.feedback.result", "Feedback received")
          .block(
            ui.section(
              ui.plain(
                "Thanks for the feedback: " +
                  ctx.values.text("summary", "value") +
                  "\nDetails: " +
                  (ctx.values.text("details", "value") || "(none)"),
              ),
            ),
          )
          .build(),
      );
    });
    command(
      "/golem-announce",
      { description: "Post a generated report file" },
      async (ctx) => {
        await ctx.slack.files.upload({
          channel_id: ctx.channelId,
          filename: "report.txt",
          content: "This report was created inside the JS bot runtime.",
        });
        return { text: "Sent report to the current channel." };
      },
    );
    action("ping.panel", async () => ({
      text: "Panel button clicked from JavaScript",
    }));
    action("ping.topic", async (ctx) => ({
      text: "Selected topic: " + ctx.action.selectedOption.value,
    }));
    action("ping.repo", async () => {});
    event("app_mention", async (ctx) => {
      await ctx.reply({ text: "I received your mention." });
    });
    event("message", async (ctx) => {
      if (ctx.text.trim() === "!pingjs")
        await ctx.reply({ text: "pong from message event" });
    });
  },
);
