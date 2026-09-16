const {defineBot} = require("slack");
module.exports = defineBot(({configure, command}) => {
  configure({name: "interaction-types", description: "Slack command text and subcommand parsing examples"});
  command("/hello", {description: "Say hello"}, async () => ({text: "Hello from the interaction types bot!"}));
  command("/echo", {description: "Echo text back"}, async ctx => ({text: ctx.text.trim().slice(0, 4000) || "Usage: /echo text"}));
  command("/fun", {description: "Try /fun roll [sides] or /fun coin"}, async ctx => {
    const args = ctx.text.trim().split(/\s+/);
    if (args[0] === "coin") return {text: Math.random() < 0.5 ? "Heads!" : "Tails!"};
    if (args[0] === "roll") {
      const sides = args[1] === undefined ? 6 : Number(args[1]);
      if (!Number.isInteger(sides) || sides < 2 || sides > 1000000) return {text: "Sides must be an integer between 2 and 1000000."};
      return {text: "You rolled " + (1 + Math.floor(Math.random() * sides)) + " (1–" + sides + ")"};
    }
    return {text: "Use /fun roll [sides] or /fun coin."};
  });
});
