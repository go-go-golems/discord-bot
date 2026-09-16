const { defineBot } = require("slack");
const ui = require("slack/ui");
module.exports = defineBot(({ configure, command, shortcut, view }) => {
  configure({
    name: "interaction-types",
    description: "Slack command text and subcommand parsing examples",
    scopes: ["users:read"],
  });
  shortcut(
    {
      callbackId: "quote-message",
      type: "message",
      name: "Quote Message",
      description: "Inspect the selected message",
    },
    async (ctx) => {
      const m = ctx.shortcut.message;
      await ctx.openModal(
        ui
          .modal("quote-result", "Quoted message")
          .block(
            ui.section(ui.plain(String(m.text || "(no text)").slice(0, 3000))),
          )
          .build(),
      );
    },
  );
  command(
    "/show-avatar",
    { description: "Select a user and view their avatar" },
    async (ctx) => {
      await ctx.openModal(
        ui
          .modal("avatar.submit", "Show avatar")
          .input("user", "User", ui.usersSelect("selected"))
          .submit("Show")
          .build(),
      );
    },
  );
  view("avatar.submit", async (ctx) => {
    const userId = ctx.values.all.user.selected.selected_user;
    const result = await ctx.slack.users.info({ user: userId });
    const user = result.user;
    const url = user.profile.image_512 || user.profile.image_192;
    await ctx.ack.update(
      ui
        .modal("avatar.result", "User avatar")
        .block(
          ui.section(
            ui.plain(user.real_name || user.name),
            ui.image(url, "User avatar"),
          ),
        )
        .build(),
    );
  });
  command("/hello", { description: "Say hello" }, async () => ({
    text: "Hello from the interaction types bot!",
  }));
  command("/echo", { description: "Echo text back" }, async (ctx) => ({
    text: ctx.text.trim().slice(0, 4000) || "Usage: /echo text",
  }));
  command(
    "/fun",
    { description: "Try /fun roll [sides] or /fun coin" },
    async (ctx) => {
      const args = ctx.text.trim().split(/\s+/);
      if (args[0] === "coin")
        return { text: Math.random() < 0.5 ? "Heads!" : "Tails!" };
      if (args[0] === "roll") {
        const sides = args[1] === undefined ? 6 : Number(args[1]);
        if (!Number.isInteger(sides) || sides < 2 || sides > 1000000)
          return { text: "Sides must be an integer between 2 and 1000000." };
        return {
          text:
            "You rolled " +
            (1 + Math.floor(Math.random() * sides)) +
            " (1–" +
            sides +
            ")",
        };
      }
      return { text: "Use /fun roll [sides] or /fun coin." };
    },
  );
});
