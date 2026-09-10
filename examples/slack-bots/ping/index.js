const { defineBot } = require("slack");

module.exports = defineBot(({ configure, command, event }) => {
  configure({
    name: "ping",
    description: "Offline Slack runtime example",
    run: { fields: {
      greeting: { type: "string", default: "pong", help: "Ping response" }
    }}
  });
  command("/golem-ping", { description: "Check the development bot" }, async (ctx) => {
    return { text: ctx.config.greeting };
  });
  event("app_mention", async (ctx) => {
    await ctx.reply({ text: "I received your mention." });
  });
});
