const {defineBot} = require("slack");

module.exports = defineBot(({configure, command}) => {
  configure({
    name: "unified-demo",
    description: "Demonstrate declared configuration and bot inspection",
    run: {fields: {
      dbPath: {type: "string", default: "./demo.sqlite", help: "Example database path; this demo does not open it"},
      apiKey: {type: "string", default: "", help: "Example external API key; never echoed"}
    }}
  });
  command("/unified-ping", {description: "Show declared runtime configuration status"}, async ctx => ({
    text: "Unified demo is alive. dbPath=" + ctx.config.dbPath +
      ", apiKey=" + (ctx.config.apiKey ? "(configured)" : "(unset)")
  }));
  command("/unified-status", {description: "Show the demo's runtime mode"}, async () => ({
    text: "active=true, mode=unified-demo, runtime=Slack Go/JavaScript"
  }));
});
