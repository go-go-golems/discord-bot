const ui = require("slack/ui"),
  data = require("./demo-store");
function article(ctx, id) {
  const original = data.getArticle(id);
  if (!original) return null;
  return Object.assign({}, original, ctx.store.get("article:" + id) || {});
}
function detail(ctx, id) {
  const a = article(ctx, id);
  if (!a) return { text: "Article not found." };
  return ui
    .message(a.title + " [" + a.status + "]")
    .block(ui.header(a.title))
    .block(ui.section(ui.plain(a.summary)))
    .block(
      ui.actions(
        "article",
        ...[
          ["verify", "Verify"],
          ["stale", "Stale"],
          ["reject", "Reject"],
          ["edit", "Edit"],
          ["source", "Source"],
          ["export", "Export"],
          ["details", "Details"],
        ].map(([name, label]) =>
          ui.button("demo.article." + name, label).value(id),
        ),
      ),
    )
    .build();
}
function search(ctx, query, page, review) {
  const matches = data
    .searchArticles(query, 50)
    .map((a) => article(ctx, a.id))
    .filter((a) => !review || a.status === review);
  const pages = Math.max(1, Math.ceil(matches.length / 3));
  page = Math.max(0, Math.min(page, pages - 1));
  const shown = matches.slice(page * 3, page * 3 + 3),
    state = { query, page, review };
  let message = ui
    .message("Articles: page " + (page + 1) + "/" + pages)
    .block(
      ui.section(
        ui.plain(
          shown.map((a) => a.title + " [" + a.status + "]").join("\n") ||
            "No matching articles",
        ),
      ),
    );
  if (shown.length)
    message.block(
      ui.actions(
        "articles",
        ui.staticSelect("demo.article.select", {
          options: shown.map((a) => ui.option(a.title, a.id)),
        }),
      ),
    );
  return message
    .block(
      ui.actions(
        "pages",
        ui
          .button("demo.search.previous", "Previous")
          .value(JSON.stringify(state)),
        ui.button("demo.search.next", "Next").value(JSON.stringify(state)),
      ),
    )
    .build();
}
function card(id) {
  const p = data.getProduct(id) || data.PRODUCTS[0];
  return ui
    .message(p.name + " · $" + p.price.toFixed(2))
    .block(ui.header(p.name))
    .block(
      ui.section(
        ui.plain(
          p.description + "\nStock: " + p.stock + " · Category: " + p.category,
        ),
      ),
    )
    .block(
      ui.actions(
        "products",
        ui.staticSelect("demo.card.select", {
          options: data.PRODUCTS.map((p) => ui.option(p.name, p.id)),
        }),
      ),
    )
    .block(
      ui.actions(
        "product-actions",
        ui.confirm(
          ui.button("demo.card.buy", "Buy").value(p.id),
          "Confirm purchase",
          "Buy " + p.name + "?",
        ),
        ui.button("demo.card.info", "Info").value(p.id),
        ui.button("demo.card.share", "Share").value(p.id),
      ),
    )
    .build();
}
function pager(page) {
  const pages = Math.ceil(data.TASKS.length / 3);
  page = Math.max(0, Math.min(page, pages - 1));
  return ui
    .message("Tasks page " + (page + 1) + "/" + pages)
    .block(
      ui.section(
        ui.plain(
          data.TASKS.slice(page * 3, page * 3 + 3)
            .map((t) => t.title + " · " + t.status)
            .join("\n"),
        ),
      ),
    )
    .block(
      ui.actions(
        "task-pages",
        ui.button("demo.pager.previous", "Previous").value(String(page)),
        ui.button("demo.pager.next", "Next").value(String(page)),
      ),
    )
    .build();
}
function register({ command, action, view, event, options }) {
  command(
    "/demo-message",
    { description: "Native message layout primitives" },
    async () =>
      ui
        .message("Native Slack message")
        .block(ui.header("Message builders"))
        .block(
          ui.section(ui.mrkdwn("*Rich text*, context, controls and URLs.")),
        )
        .block(
          ui.context(
            ui.plain("Neutral, primary and danger are Slack's button styles."),
          ),
        )
        .block(
          ui.actions(
            "buttons",
            ui.button("demo.primary", "Primary").style("primary"),
            ui.button("demo.neutral", "Neutral"),
            ui.button("demo.danger", "Danger").style("danger"),
            ui.linkButton(
              "demo.link",
              "Source",
              "https://github.com/go-go-golems/discord-bot",
            ),
          ),
        )
        .build(),
  );
  for (const name of ["primary", "neutral", "danger"])
    action("demo." + name, async () => ({ text: "Clicked " + name }));
  action("demo.link", async () => {});
  command(
    "/demo-form",
    { description: "Open a native modal form" },
    async (ctx) => {
      await ctx.openModal(
        ui
          .modal("demo.form", "Example form")
          .input("title", "Title", ui.textInput("value").length(3, 100))
          .input("body", "Body", ui.textInput("value").multiline(), {
            optional: true,
          })
          .submit("Save")
          .build(),
      );
    },
  );
  view("demo.form", async (ctx) => {
    const title = ctx.values.text("title", "value") || "";
    if (title.trim().length < 3) {
      await ctx.ack.errors({ title: "Use three or more characters." });
      return;
    }
    await ctx.ack.update(
      ui
        .modal("demo.form.result", "Submitted")
        .block(
          ui.section(
            ui.plain(title + "\n" + (ctx.values.text("body", "value") || "")),
          ),
        )
        .build(),
    );
  });
  for (const name of ["demo-search", "find"])
    command("/" + name, { description: "Search articles" }, async (ctx) =>
      search(ctx, ctx.text.trim(), 0, ""),
    );
  command(
    "/demo-review",
    { description: "Review articles [status]" },
    async (ctx) => search(ctx, "", 0, ctx.text.trim() || "draft"),
  );
  for (const direction of ["previous", "next"])
    action("demo.search." + direction, async (ctx) => {
      const s = JSON.parse(ctx.action.value);
      await ctx.replaceOriginal(
        search(
          ctx,
          s.query,
          s.page + (direction === "next" ? 1 : -1),
          s.review,
        ),
      );
    });
  action("demo.article.select", async (ctx) =>
    detail(ctx, ctx.action.selectedOption.value),
  );
  for (const [name, status] of [
    ["verify", "verified"],
    ["stale", "stale"],
    ["reject", "rejected"],
  ])
    action("demo.article." + name, async (ctx) => {
      const a = article(ctx, ctx.action.value);
      a.status = status;
      ctx.store.set("article:" + a.id, a);
      await ctx.replaceOriginal(detail(ctx, a.id));
    });
  action("demo.article.edit", async (ctx) => {
    const a = article(ctx, ctx.action.value);
    await ctx.openModal(
      ui
        .modal("demo.article.save", "Edit article")
        .metadata(a.id)
        .input("title", "Title", ui.textInput("value").initial(a.title))
        .input(
          "summary",
          "Summary",
          ui.textInput("value").multiline().initial(a.summary),
        )
        .submit("Save")
        .build(),
    );
  });
  view("demo.article.save", async (ctx) => {
    const a = article(ctx, ctx.view.privateMetadata);
    a.title = ctx.values.text("title", "value");
    a.summary = ctx.values.text("summary", "value");
    if (!a.title.trim()) {
      await ctx.ack.errors({ title: "Title required." });
      return;
    }
    ctx.store.set("article:" + a.id, a);
    await ctx.ack.update(
      ui
        .modal("demo.article.saved", "Saved")
        .block(ui.section(ui.plain(a.title)))
        .build(),
    );
  });
  action("demo.article.source", async () => ({
    text: "Demonstration data: examples/slack-bots/ui-showcase/lib/demo-store.js",
  }));
  action("demo.article.details", async (ctx) => ({
    text: JSON.stringify(article(ctx, ctx.action.value), null, 2),
  }));
  action("demo.article.export", async (ctx) => {
    const a = article(ctx, ctx.action.value);
    await ctx.slack.files.upload({
      channel_id: ctx.channelId,
      filename: a.id + ".md",
      content: "# " + a.title + "\n\n" + a.summary,
    });
    return { text: "Exported " + a.title };
  });
  command(
    "/demo-confirm",
    { description: "Inline yes/no confirmation" },
    async () =>
      ui
        .message("Confirm this demo action?")
        .block(
          ui.actions(
            "confirm",
            ui.button("demo.confirm.yes", "Confirm").style("danger"),
            ui.button("demo.confirm.no", "Cancel"),
          ),
        )
        .build(),
  );
  for (const [name, text] of [
    ["yes", "Confirmed."],
    ["no", "Cancelled."],
  ])
    action("demo.confirm." + name, async (ctx) => {
      await ctx.replaceOriginal({ text });
    });
  command("/demo-pager", { description: "Paginated task list" }, async () =>
    pager(0),
  );
  for (const direction of ["previous", "next"])
    action("demo.pager." + direction, async (ctx) => {
      await ctx.replaceOriginal(
        pager(Number(ctx.action.value) + (direction === "next" ? 1 : -1)),
      );
    });
  for (const name of ["demo-cards", "browse"])
    command("/" + name, { description: "Browse product cards" }, async () =>
      card(data.PRODUCTS[0].id),
    );
  action("demo.card.select", async (ctx) => {
    await ctx.replaceOriginal(card(ctx.action.selectedOption.value));
  });
  action("demo.card.buy", async (ctx) => {
    const p = data.getProduct(ctx.action.value);
    return {
      text: p.stock
        ? "Demo purchase confirmed: " + p.name
        : "Out of stock: " + p.name,
    };
  });
  action("demo.card.info", async (ctx) => ({
    text: JSON.stringify(data.getProduct(ctx.action.value), null, 2),
  }));
  action("demo.card.share", async (ctx) => {
    await ctx.slack.messages.post(
      Object.assign({ channelId: ctx.channelId }, card(ctx.action.value)),
    );
    return { text: "Shared product in this channel." };
  });
  command(
    "/demo-selects",
    { description: "Static, user, group, channel and conversation selections" },
    async () =>
      ui
        .message("Native selections")
        .block(
          ui.actions(
            "static",
            ui.staticSelect("demo.select.string", {
              options: [ui.option("One", "1"), ui.option("Two", "2")],
            }),
          ),
        )
        .block(
          ui.actions(
            "users",
            ui.usersSelect("demo.select.user"),
            ui.multiUsersSelect("demo.select.users"),
          ),
        )
        .block(
          ui.actions(
            "channels",
            ui.channelsSelect("demo.select.channel"),
            ui.conversationsSelect("demo.select.conversation"),
          ),
        )
        .block(
          ui.actions(
            "groups",
            ui.externalSelect("demo.select.group", {
              placeholder: "User group",
              min_query_length: 0,
            }),
          ),
        )
        .build(),
  );
  options("demo.select.group", async (ctx) => {
    const result = await ctx.slack.usergroups.list({});
    await ctx.ack.options(
      (result.usergroups || [])
        .filter((g) =>
          g.name.toLowerCase().includes((ctx.query || "").toLowerCase()),
        )
        .slice(0, 100)
        .map((g) => ui.option(g.name, g.id)),
    );
  });
  for (const kind of [
    "string",
    "user",
    "users",
    "channel",
    "conversation",
    "group",
  ])
    action("demo.select." + kind, async (ctx) => ({
      text: "Selection: " + JSON.stringify(ctx.action.selection),
    }));
  const alias = async () => ({
    text: "Both command names invoke the same handler.",
  });
  command("/demo-alias", { description: "Shared handler example" }, alias);
  command(
    "/demo-alias-alt",
    { description: "Same handler under another command" },
    alias,
  );
  event("message", async (ctx) => {
    if (ctx.text.trim() === "!showcase")
      await ctx.reply({
        text: "Try /demo-message, /demo-form, /demo-search, /demo-review, /demo-confirm, /demo-pager, /demo-cards or /demo-selects.",
      });
  });
}
module.exports = { register };
