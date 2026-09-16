const { defineBot } = require("slack"),
  ui = require("slack/ui"),
  { createStore } = require("../lib/sqlite");
const store = createStore([
  "CREATE TABLE IF NOT EXISTS links (team TEXT NOT NULL,id TEXT NOT NULL,url TEXT NOT NULL,title TEXT NOT NULL,summary TEXT NOT NULL,tags TEXT NOT NULL,author TEXT NOT NULL,updated TEXT NOT NULL,PRIMARY KEY(team,id),UNIQUE(team,url))",
]);
function find(ctx, id) {
  store.ensure(ctx);
  return store.query(
    "SELECT * FROM links WHERE team=? AND id=?",
    ctx.teamId,
    id,
  )[0];
}
function search(ctx, q) {
  store.ensure(ctx);
  const like = "%" + q + "%";
  return store.query(
    "SELECT * FROM links WHERE team=? AND (title LIKE ? OR summary LIKE ? OR tags LIKE ? OR url LIKE ?) ORDER BY updated DESC LIMIT ?",
    ctx.teamId,
    like,
    like,
    like,
    like,
    Math.max(1, Math.min(100, ctx.config.searchLimit)),
  );
}
function save(ctx, values) {
  store.ensure(ctx);
  const [url, title, summary, tags] = values.map((v) => String(v || "").trim());
  if (!/^https:\/\//.test(url) || title.length < 2)
    throw Error("Use an HTTPS URL and a title of at least two characters.");
  const existing = store.query(
    "SELECT id FROM links WHERE team=? AND url=?",
    ctx.teamId,
    url,
  )[0];
  const id = existing
    ? existing.id
    : Date.now().toString(36) + Math.random().toString(36).slice(2, 7);
  store.exec(
    "INSERT INTO links VALUES (?,?,?,?,?,?,?,?) ON CONFLICT(team,url) DO UPDATE SET title=excluded.title,summary=excluded.summary,tags=excluded.tags,author=excluded.author,updated=excluded.updated",
    ctx.teamId,
    id,
    url,
    title,
    summary,
    tags,
    ctx.userId,
    new Date().toISOString(),
  );
  return find(ctx, id);
}
function detail(link) {
  if (!link) return { text: "Link not found." };
  return ui
    .message("Saved link: " + link.title)
    .block(ui.header(link.title.slice(0, 150)))
    .block(ui.section(ui.plain((link.summary || link.url).slice(0, 3000))))
    .block(
      ui.context(
        ui.plain("Tags: " + (link.tags || "none") + " · ID: " + link.id),
      ),
    )
    .block(
      ui.actions("source", ui.linkButton("kb.open", "Open link", link.url)),
    )
    .build();
}
function listing(ctx, q) {
  ctx.store.set("custom-kb.query:" + ctx.userId, q);
  const links = search(ctx, q);
  let message = ui
    .message(
      links.length
        ? "Found " + links.length + " link(s)."
        : "No links found. Use /kb-add or /kb-link.",
    )
    .block(
      ui.section(
        ui.plain(
          (
            links.map((l) => l.title + " — " + l.url).join("\n") || "No results"
          ).slice(0, 3000),
        ),
      ),
    );
  if (links.length)
    message.block(
      ui.actions(
        "results",
        ui.staticSelect("kb.select", {
          options: links.map((l) => ui.option(l.title.slice(0, 75), l.id)),
        }),
      ),
    );
  return message
    .block(ui.actions("refresh", ui.button("kb.refresh", "Refresh")))
    .build();
}
module.exports = defineBot(({ configure, command, action, view, options }) => {
  configure({
    name: "custom-kb",
    description: "Persistent workspace-scoped link collection",
    run: {
      fields: {
        dbPath: { type: "string", default: "./slack-custom-kb.sqlite" },
        searchLimit: { type: "number", default: 10 },
      },
    },
  });
  command("/kb-add", { description: "Add a link with a form" }, async (ctx) => {
    await ctx.openModal(
      ui
        .modal("kb.add", "Add KB link")
        .input("url", "HTTPS URL", ui.textInput("value"))
        .input("title", "Title", ui.textInput("value").length(2, 120))
        .input("summary", "Summary", ui.textInput("value").multiline(), {
          optional: true,
        })
        .input("tags", "Tags", ui.textInput("value"), { optional: true })
        .submit("Save")
        .build(),
    );
  });
  view("kb.add", async (ctx) => {
    let link;
    try {
      link = save(
        ctx,
        ["url", "title", "summary", "tags"].map((id) =>
          ctx.values.text(id, "value"),
        ),
      );
    } catch (err) {
      await ctx.ack.errors({ url: err.message });
      return;
    }
    await ctx.ack.update(
      ui
        .modal("kb.saved", "Link saved")
        .block(ui.section(ui.plain(link.title + "\n" + link.url)))
        .build(),
    );
  });
  command(
    "/kb-link",
    { description: "url | title | summary | tags" },
    async (ctx) => {
      try {
        return detail(save(ctx, ctx.text.split("|")));
      } catch (err) {
        return { text: err.message };
      }
    },
  );
  command("/kb-list", { description: "Recent links" }, async (ctx) =>
    listing(ctx, ""),
  );
  command(
    "/kb-search",
    { description: "Search links; omit text for live suggestions" },
    async (ctx) => {
      if (ctx.text.trim()) return listing(ctx, ctx.text.trim());
      await ctx.openModal(
        ui
          .modal("kb.search", "Search links")
          .input(
            "link",
            "Link",
            ui.externalSelect("kb.suggest", { min_query_length: 0 }),
          )
          .submit("Open")
          .build(),
      );
    },
  );
  options("kb.suggest", async (ctx) => {
    await ctx.ack.options(
      search(ctx, ctx.query || "").map((l) =>
        ui.option(l.title.slice(0, 75), l.id),
      ),
    );
  });
  view("kb.search", async (ctx) => {
    const selected = ctx.values.all.link["kb.suggest"].selected_option;
    const link = find(ctx, selected.value);
    await ctx.ack.update(
      ui
        .modal("kb.found", "Stored link")
        .block(
          ui.section(
            ui.plain(link ? link.title + "\n" + link.url : "Link not found"),
          ),
        )
        .build(),
    );
  });
  action("kb.select", async (ctx) =>
    detail(find(ctx, ctx.action.selectedOption.value)),
  );
  action("kb.refresh", async (ctx) => {
    await ctx.replaceOriginal(
      listing(ctx, ctx.store.get("custom-kb.query:" + ctx.userId) || ""),
    );
  });
  action("kb.open", async () => {});
});
