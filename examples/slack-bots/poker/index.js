const { defineBot } = require("slack"),
  ui = require("slack/ui"),
  poker = require("./lib/poker");
const help =
  "/poker-deal · /poker-draw 1,3,5 · /poker-score · /poker-reset\n/poker-rank As Kd Qh Jc Tc\n/poker-action As Kd | Qh Jc Tc | 120 | 30 (hole | board | pot | call)";
function panel(text) {
  return ui
    .message(text)
    .block(ui.section(ui.plain(text)))
    .block(
      ui.actions(
        "poker",
        ...[
          ["deal", "Deal"],
          ["score", "Score"],
          ["rank", "Rank"],
          ["action", "Advice"],
          ["reset", "Reset"],
        ].map(([id, label]) => ui.button("poker." + id, label)),
      ),
    )
    .build();
}
function rank(text) {
  const cards = poker.parseCardList(text);
  if (cards.length < 5 || cards.length > 7) throw Error("Use 5–7 cards.");
  poker.validateNoDuplicateCards(cards);
  const best = poker.chooseBestHand(cards);
  return poker.renderHandSummary(best.cards, best);
}
function advice(hole, board, pot, toCall) {
  const result = poker.suggestAction({ hole, board, pot, toCall });
  return poker.renderAdvice(result.hole, result.board, result);
}
function round(ctx, op) {
  const state = poker.pokerState(ctx);
  if (op === "reset") {
    poker.resetRound(state);
    return "Round cleared.";
  }
  if (op === "deal") poker.startRound(state);
  if (op === "draw") poker.redrawRound(state, poker.parsePositions(ctx.text));
  const summary = poker.stateSummary(state);
  return summary
    ? poker.renderHandSummary(summary.hand, summary.best) +
        "\nRound " +
        summary.round +
        " · draw used: " +
        summary.drawn
    : "No hand. Use /poker-deal.";
}
const safe = (fn) => async (ctx) => {
  try {
    return panel(fn(ctx));
  } catch (err) {
    return { text: err.message };
  }
};
function form(kind) {
  let modal = ui.modal(
    "poker." + kind + ".submit",
    kind === "rank" ? "Rank hand" : "Action advice",
  );
  if (kind === "rank")
    modal.input("cards", "Five to seven cards", ui.textInput("value"));
  else
    modal
      .input("hole", "Two hole cards", ui.textInput("value"))
      .input("board", "Board cards", ui.textInput("value"), { optional: true })
      .input("pot", "Pot", ui.textInput("value").initial("0"))
      .input("call", "To call", ui.textInput("value").initial("0"));
  return modal.submit("Evaluate").build();
}
module.exports = defineBot(({ configure, command, action, view, event }) => {
  configure({
    name: "poker",
    description: "Five-card draw, hand ranking and Hold'em advice",
  });
  command(
    "/poker-help",
    { description: "Poker commands and quick controls" },
    async () => panel(help),
  );
  for (const op of ["deal", "draw", "score", "reset"])
    command(
      "/poker-" + op,
      { description: "Poker " + op },
      safe((ctx) => round(ctx, op)),
    );
  command(
    "/poker-rank",
    { description: "Rank 5–7 cards" },
    safe((ctx) => rank(ctx.text)),
  );
  command(
    "/poker-action",
    { description: "hole | board | pot | to-call" },
    safe((ctx) => advice(...ctx.text.split("|"))),
  );
  for (const op of ["deal", "score", "reset"])
    action(
      "poker." + op,
      safe((ctx) => round(ctx, op)),
    );
  for (const kind of ["rank", "action"]) {
    action("poker." + kind, async (ctx) => {
      await ctx.openModal(form(kind));
    });
    view("poker." + kind + ".submit", async (ctx) => {
      const value = (id) => ctx.values.text(id, "value") || "";
      let text;
      try {
        text =
          kind === "rank"
            ? rank(value("cards"))
            : advice(
                value("hole"),
                value("board"),
                value("pot"),
                value("call"),
              );
      } catch (err) {
        await ctx.ack.errors({
          [kind === "rank" ? "cards" : "hole"]: err.message,
        });
        return;
      }
      await ctx.ack.update(
        ui
          .modal("poker.result", "Poker result")
          .block(ui.section(ui.plain(text)))
          .close("Close")
          .build(),
      );
    });
  }
  event("message", async (ctx) => {
    if (ctx.text.trim() === "!poker") await ctx.reply({ text: help });
  });
});
