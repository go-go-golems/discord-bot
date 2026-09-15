const { defineBot } = require("slack");
const ui = require("slack/ui");

module.exports = defineBot(({ configure, command, event, action, view }) => {
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

  action("showcase.ack", async ctx => {
    await ctx.openModal(
      ui.modal("showcase.edit", "Edit showcase")
        .metadata("showcase")
        .input("title", "Title", ui.textInput("title_input").initial("Acknowledge"))
        .submit("Save")
        .build(),
    );
  });

  view("showcase.edit", async ctx => {
    const title = ctx.values.text("title", "title_input");
    if (!title || title.trim().length < 3) {
      return ctx.ack.errors({ title: "Use at least three characters." });
    }
    await ctx.ack.accept();
    ctx.log.info(`Saved ${title.trim()}`);
  });

  event("app_mention", async ctx => ({
    text: "The UI showcase is available through /ui-showcase.",
  }));
});
