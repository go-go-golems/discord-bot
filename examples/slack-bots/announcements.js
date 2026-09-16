const { defineBot } = require("slack");
const ui = require("slack/ui");

module.exports = defineBot(({configure, command}) => {
  configure({name: "announcements", description: "Preview announcements with Slack Block Kit"});
  command("/announce-preview", {description: "Preview an announcement: /announce-preview title"}, async ctx => {
    const title = ctx.text.trim();
    if (!title) return {text: "Usage: /announce-preview title"};
    return ui.message("Preview ready for " + title.slice(0, 150))
      .block(ui.header(title.slice(0, 150)))
      .block(ui.section(ui.plain("Announcement preview generated from a root-level bot script.")))
      .build();
  });
});
