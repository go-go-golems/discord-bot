const { defineBot } = require("slack");
const ui = require("slack/ui");

module.exports = defineBot(({ configure, command, event }) => {
  configure({
    name: "ui-showcase",
    description: "Demonstrate Slack Block Kit messages",
  });

  command("/ui-showcase", { description: "Post a Block Kit message" }, async ctx => {
    const message = ui
      .message("UI showcase: choose an action")
      .block(ui.header("Slack UI showcase"))
      .block(ui.section(ui.mrkdwn("This message was built by the native Slack UI module.")))
      .block(ui.actions(
        "showcase-actions",
        ui.button("showcase.ack", "Acknowledge").value("showcase"),
      ))
      .build();

    return message;
  });

  event("app_mention", async ctx => ({
    text: "The UI showcase is available through /ui-showcase.",
  }));
});
