const {defineBot} = require("slack");
const ui = require("slack/ui");
const ROASTS = [
  "I would explain it to you, but I left my crayons in another workspace.",
  "Your aura has the latency of a potato on hotel Wi-Fi.",
  "I checked the logs. Even the warnings are disappointed.",
  "You're not the problem. You're the regression test for the problem.",
  "I have seen better decisions from a random number generator.",
  "If confidence compiled code, you'd still have syntax errors."
];
const BACKHANDED = [
  "You are technically present, which is more than I expected.",
  "That was almost a thought. Proud of you.",
  "You're doing your best, and statistics says that has to be enough sometimes.",
  "I admire your commitment to making me lower my expectations."
];
const pick = items => items[Math.floor(Math.random() * items.length)];
const roast = () => ({text: pick(ROASTS)});
module.exports = defineBot(({configure, command, action, view, event}) => {
  configure({name: "hater", description: "Playful self-roasts, reluctant compliments and apology forms"});
  command("/hate", {description: "Receive the bot's current level of contempt"}, async () =>
    ui.message("Contempt report")
      .block(ui.header("Contempt report"))
      .block(ui.section(ui.plain("I hate you " + Math.floor(Math.random() * 101) + "%. " + pick(ROASTS))))
      .block(ui.actions("hater-actions",
        ui.button("hater.roast", "Roast me again").style("danger"),
        ui.button("hater.mercy", "Beg for mercy"),
        ui.button("hater.apology", "Write apology")))
      .build());
  command("/roast", {description: "Ask the bot to roast you"}, async () => roast());
  command("/compliment", {description: "Receive a reluctant compliment"}, async () => ({text: pick(BACKHANDED)}));
  action("hater.roast", async () => roast());
  action("hater.mercy", async () => ({text: "Mercy denied. Try having fewer opinions next time."}));
  action("hater.apology", async ctx => {
    await ctx.openModal(ui.modal("hater.apology.submit", "Write an apology")
      .input("subject", "What are you apologizing for?", ui.textInput("value"))
      .input("details", "Make it convincing", ui.textInput("value"))
      .submit("Submit").build());
  });
  view("hater.apology.submit", async ctx => {
    const subject = (ctx.values.text("subject", "value") || "").trim();
    const details = (ctx.values.text("details", "value") || "").trim();
    const errors = {};
    if (subject.length < 3) errors.subject = "Use at least three characters.";
    if (details.length < 10) errors.details = "Make the apology at least ten characters.";
    if (Object.keys(errors).length) return ctx.ack.errors(errors);
    ctx.store.set("apology:" + ctx.userId, {subject, details});
    await ctx.ack.accept();
  });
  command("/apology-status", {description: "Read the verdict on your latest apology"}, async ctx => ({
    text: ctx.store.get("apology:" + ctx.userId)
      ? "Apology received. Forgiveness not found."
      : "No apology on record. Open /hate and choose Write apology."
  }));
  event("app_mention", async () => roast());
});
