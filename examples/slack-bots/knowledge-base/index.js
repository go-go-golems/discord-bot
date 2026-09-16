const { defineBot } = require("slack");
const ui = require("slack/ui");
const { createStore } = require("../lib/sqlite");
const { authorized, ids } = require("../lib/authorized");
const capture = require("./lib/capture");
const store = createStore([
  "CREATE TABLE IF NOT EXISTS knowledge (team TEXT,id TEXT,title TEXT,summary TEXT,body TEXT,tags TEXT,aliases TEXT,status TEXT,source TEXT,channel TEXT,ts TEXT,author TEXT,updated TEXT,PRIMARY KEY(team,id))",
  "CREATE INDEX IF NOT EXISTS knowledge_source ON knowledge(team,channel,ts)",
  "CREATE TABLE IF NOT EXISTS knowledge_audit(team TEXT,id TEXT,actor TEXT,action TEXT,created TEXT)",
]);
function ensure(ctx) {
  store.ensure(ctx);
  if (
    ctx.config.seedEntries &&
    !store.query("SELECT id FROM knowledge WHERE team=? LIMIT 1", ctx.teamId)
      .length
  ) {
    save(ctx, {
      id: "getting-started",
      title: "Getting started",
      summary: "Use /teach to save knowledge.",
      body: "Use /teach to record knowledge, /ask to search it, and /review to verify candidates.",
      status: "verified",
    });
  }
}
function get(ctx, id) {
  ensure(ctx);
  const exact = store.query(
    "SELECT * FROM knowledge WHERE team=? AND id=?",
    ctx.teamId,
    id,
  )[0];
  if (exact) return exact;
  const name = String(id || "")
    .trim()
    .toLowerCase();
  return store.query("SELECT * FROM knowledge WHERE team=?", ctx.teamId).find(
    (entry) =>
      entry.title.toLowerCase() === name ||
      entry.title
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/^-|-$/g, "") === name ||
      String(entry.aliases || "")
        .toLowerCase()
        .split(",")
        .map((s) => s.trim())
        .includes(name),
  );
}

function save(ctx, entry) {
  store.ensure(ctx);
  const id =
    entry.id ||
    Date.now().toString(36) + Math.random().toString(36).slice(2, 7);
  store.exec(
    "INSERT INTO knowledge VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(team,id) DO UPDATE SET title=excluded.title,summary=excluded.summary,body=excluded.body,tags=excluded.tags,aliases=excluded.aliases,status=excluded.status,source=excluded.source,updated=excluded.updated",
    ctx.teamId,
    id,
    entry.title,
    entry.summary || "",
    entry.body,
    entry.tags || "",
    entry.aliases || "",
    entry.status || "draft",
    entry.source || "",
    entry.channel || ctx.channelId || "",
    entry.ts || "",
    ctx.userId,
    new Date().toISOString(),
  );
  audit(ctx, id, "save");
  return store.query(
    "SELECT * FROM knowledge WHERE team=? AND id=?",
    ctx.teamId,
    id,
  )[0];
}
function audit(ctx, id, action) {
  store.exec(
    "INSERT INTO knowledge_audit VALUES (?,?,?,?,?)",
    ctx.teamId,
    id,
    ctx.userId,
    action,
    new Date().toISOString(),
  );
}
async function reviewer(ctx) {
  return authorized(
    ctx,
    ctx.config.trustedReviewerIds,
    ctx.config.trustedReviewerGroups,
    true,
  );
}
async function status(ctx, id, value) {
  if (!(await reviewer(ctx)))
    throw Error("You are not an authorized reviewer.");
  const entry = get(ctx, id);
  if (!entry) throw Error("Knowledge entry not found.");
  store.exec(
    "UPDATE knowledge SET status=?,updated=? WHERE team=? AND id=?",
    value,
    new Date().toISOString(),
    ctx.teamId,
    id,
  );
  audit(ctx, id, value);
  return get(ctx, id);
}
function detail(entry) {
  if (!entry) return { text: "Knowledge entry not found." };
  return ui
    .message(entry.title + " [" + entry.status + "]")
    .block(ui.header(entry.title.slice(0, 150)))
    .block(ui.section(ui.plain(entry.body.slice(0, 3000))))
    .block(
      ui.context(
        ui.plain("ID: " + entry.id + " · " + entry.status + " · " + entry.tags),
      ),
    )
    .block(
      ui.actions(
        "review",
        ...[
          ["verify", "Verify"],
          ["stale", "Stale"],
          ["reject", "Reject"],
          ["edit", "Edit"],
          ["source", "Source"],
          ["export", "Export"],
        ].map(([action, label]) =>
          ui.button("knowledge." + action, label).value(entry.id),
        ),
      ),
    )
    .build();
}
function rows(ctx, query, filter) {
  ensure(ctx);
  const pattern = "%" + query + "%";
  return store.query(
    "SELECT * FROM knowledge WHERE team=? AND (title LIKE ? OR summary LIKE ? OR body LIKE ? OR tags LIKE ? OR aliases LIKE ?) AND (?='' OR status=?) ORDER BY updated DESC LIMIT 1000",
    ctx.teamId,
    pattern,
    pattern,
    pattern,
    pattern,
    pattern,
    filter,
    filter,
  );
}
function listing(ctx, query, filter, page) {
  const all = rows(ctx, query, filter),
    size = Math.max(1, Math.min(25, Math.floor(ctx.config.reviewLimit)));
  page = Math.max(
    0,
    Math.min(page, Math.max(0, Math.ceil(all.length / size) - 1)),
  );
  const entries = all.slice(page * size, (page + 1) * size);
  const msg = ui
    .message("Knowledge: " + all.length + " matches · page " + (page + 1))
    .block(
      ui.section(
        ui.plain(
          entries
            .map((e) => e.id + " · " + e.title + " [" + e.status + "]")
            .join("\n") || "No matching entries.",
        ),
      ),
    );
  if (entries.length)
    msg.block(
      ui.actions(
        "select",
        ui.staticSelect("knowledge.select", {
          options: entries.map((e) => ui.option(e.title.slice(0, 75), e.id)),
        }),
      ),
    );
  msg.block(
    ui.actions(
      "pages",
      ui
        .button("knowledge.previous", "Previous")
        .value(JSON.stringify({ query, filter, page })),
      ui
        .button("knowledge.next", "Next")
        .value(JSON.stringify({ query, filter, page })),
    ),
  );
  return msg.build();
}
function form(entry) {
  const view = ui
    .modal("knowledge.submit", entry ? "Edit knowledge" : "Teach knowledge")
    .metadata(entry ? entry.id : "");
  for (const [field, label, min, max] of [
    ["title", "Title", 3, 100],
    ["summary", "Summary", 10, 300],
    ["body", "Body", 20, 2000],
    ["tags", "Tags", 0, 200],
    ["aliases", "Aliases", 0, 200],
    ["source", "Source URL or note", 0, 300],
  ]) {
    let input = ui
      .textInput("value")
      .length(min, max)
      .initial(entry ? entry[field] || "" : "");
    if (field === "body" || field === "summary") input.multiline();
    view.input(field, label, input, { optional: min === 0 });
  }
  return view.submit("Save").build();
}
module.exports = defineBot(
  ({ configure, command, action, view, event, options }) => {
    configure({
      name: "knowledge-base",
      description: "Capture, search and curate persistent Slack knowledge",
      scopes: ["files:write", "usergroups:read"],
      run: {
        fields: {
          dbPath: { type: "string", default: "./slack-knowledge.sqlite" },
          captureEnabled: { type: "bool", default: true },
          captureThreshold: { type: "number", default: 0.65 },
          captureChannels: { type: "string", default: "" },
          reviewLimit: { type: "number", default: 5 },
          seedEntries: { type: "bool", default: true },
          reactionPromoteEmojis: { type: "string", default: "brain,pushpin" },
          trustedReviewerIds: { type: "string", default: "" },
          trustedReviewerGroups: { type: "string", default: "" },
        },
      },
    });
    for (const name of ["teach", "remember"])
      command(
        "/" + name,
        { description: "Teach the shared knowledge base" },
        async (ctx) => {
          await ctx.openModal(form());
        },
      );
    view("knowledge.submit", async (ctx) => {
      const id = ctx.view.privateMetadata;
      if (id && !(await reviewer(ctx))) {
        await ctx.ack.errors({ title: "Reviewer permission required." });
        return;
      }
      const entry = id ? get(ctx, id) : {};
      if (!entry) {
        await ctx.ack.errors({ title: "Entry no longer exists." });
        return;
      }
      for (const key of [
        "title",
        "summary",
        "body",
        "tags",
        "aliases",
        "source",
      ])
        entry[key] = (ctx.values.text(key, "value") || "").trim();
      const errors = {};
      if (entry.title.length < 3)
        errors.title = "Use at least three characters.";
      if (entry.summary.length < 10)
        errors.summary = "Use at least ten characters.";
      if (entry.body.length < 20)
        errors.body = "Use at least twenty characters.";
      if (Object.keys(errors).length) {
        await ctx.ack.errors(errors);
        return;
      }
      const saved = save(ctx, entry);
      await ctx.ack.update(
        ui
          .modal("knowledge.saved", "Knowledge saved")
          .block(ui.section(ui.plain(saved.title + "\nID: " + saved.id)))
          .build(),
      );
    });
    for (const name of ["ask", "kb-search"])
      command(
        "/" + name,
        { description: "Search knowledge; empty text opens suggestions" },
        async (ctx) => {
          if (ctx.text.trim()) return listing(ctx, ctx.text.trim(), "", 0);
          await ctx.openModal(
            ui
              .modal("knowledge.open", "Find knowledge")
              .input(
                "entry",
                "Entry",
                ui.externalSelect("knowledge.suggest", { min_query_length: 0 }),
              )
              .submit("Open")
              .build(),
          );
        },
      );
    options("knowledge.suggest", async (ctx) => {
      await ctx.ack.options(
        rows(ctx, ctx.query || "", "")
          .slice(0, 100)
          .map((e) => ui.option(e.title.slice(0, 75), e.id)),
      );
    });
    view("knowledge.open", async (ctx) => {
      const entry = get(
        ctx,
        ctx.values.all.entry["knowledge.suggest"].selected_option.value,
      );
      await ctx.ack.update(
        ui
          .modal("knowledge.result", "Knowledge")
          .block(
            ui.section(
              ui.plain(entry ? entry.body.slice(0, 3000) : "Not found"),
            ),
          )
          .build(),
      );
    });
    for (const name of ["article", "kb-article"])
      command(
        "/" + name,
        { description: "Read an entry by ID or slug; empty opens suggestions" },
        async (ctx) => {
          if (ctx.text.trim()) return detail(get(ctx, ctx.text.trim()));
          await ctx.openModal(
            ui
              .modal("knowledge.open", "Find knowledge")
              .input(
                "entry",
                "Entry",
                ui.externalSelect("knowledge.suggest", { min_query_length: 0 }),
              )
              .submit("Open")
              .build(),
          );
        },
      );
    for (const name of ["review", "kb-review"])
      command(
        "/" + name,
        { description: "Review entries by status (default draft)" },
        async (ctx) => listing(ctx, "", ctx.text.trim() || "draft", 0),
      );
    for (const name of ["recent", "kb-recent"])
      command("/" + name, { description: "Recent knowledge" }, async (ctx) =>
        listing(ctx, "", "", 0),
      );
    for (const [name, next] of [
      ["verify", "verified"],
      ["stale", "stale"],
      ["reject", "rejected"],
    ]) {
      command(
        "/kb-" + name,
        { description: name + " an entry by ID" },
        async (ctx) => {
          try {
            return detail(await status(ctx, ctx.text.trim(), next));
          } catch (err) {
            return { text: err.message };
          }
        },
      );
      action("knowledge." + name, async (ctx) => {
        try {
          await ctx.replaceOriginal(
            detail(await status(ctx, ctx.action.value, next)),
          );
        } catch (err) {
          await ctx.reply({ text: err.message });
        }
      });
    }
    action("knowledge.select", async (ctx) =>
      detail(get(ctx, ctx.action.selectedOption.value)),
    );
    for (const direction of ["previous", "next"])
      action("knowledge." + direction, async (ctx) => {
        const state = JSON.parse(ctx.action.value);
        await ctx.replaceOriginal(
          listing(
            ctx,
            state.query,
            state.filter,
            state.page + (direction === "next" ? 1 : -1),
          ),
        );
      });
    action("knowledge.edit", async (ctx) => {
      if (!(await reviewer(ctx)))
        return { text: "Reviewer permission required." };
      const entry = get(ctx, ctx.action.value);
      if (!entry) return { text: "Entry not found." };
      await ctx.openModal(form(entry));
    });
    action("knowledge.source", async (ctx) => {
      const entry = get(ctx, ctx.action.value);
      return {
        text: entry
          ? entry.source || "No source recorded."
          : "Entry not found.",
      };
    });
    action("knowledge.export", async (ctx) => {
      const entry = get(ctx, ctx.action.value);
      if (!entry) return { text: "Entry not found." };
      await ctx.slack.files.upload({
        channel_id: ctx.channelId,
        filename: entry.id + ".md",
        content:
          "# " +
          entry.title +
          "\n\n" +
          entry.body +
          "\n\nSource: " +
          entry.source,
      });
      return { text: "Exported " + entry.id };
    });
    event("message", async (ctx) => {
      if (
        !ctx.config.captureEnabled ||
        (ids(ctx.config.captureChannels).length &&
          !ids(ctx.config.captureChannels).includes(ctx.channelId))
      )
        return;
      if (capture.scoreMessage(ctx.text) < ctx.config.captureThreshold) return;
      ensure(ctx);
      if (
        store.query(
          "SELECT id FROM knowledge WHERE team=? AND channel=? AND ts=?",
          ctx.teamId,
          ctx.channelId,
          ctx.event.ts,
        ).length
      )
        return;
      const link = await ctx.slack.messages.permalink({
        conversation_id: ctx.channelId,
        message_ts: ctx.event.ts,
      });
      const entry = save(ctx, {
        title: capture.deriveTitle(ctx.text),
        summary: capture.deriveSummary(ctx.text),
        body: ctx.text,
        tags: capture.deriveTags(ctx.text).join(","),
        aliases: capture
          .deriveAliases(ctx.text, capture.deriveTitle(ctx.text))
          .join(","),
        channel: ctx.channelId,
        ts: ctx.event.ts,
        source: link.permalink || "",
      });
      await ctx.reply(detail(entry));
    });
    event("reaction_added", async (ctx) => {
      const data = ctx.event.data;
      if (
        !ids(ctx.config.reactionPromoteEmojis).includes(data.reaction) ||
        !(await reviewer(ctx)) ||
        !data.item
      )
        return;
      ensure(ctx);
      const entry = store.query(
        "SELECT * FROM knowledge WHERE team=? AND channel=? AND ts=?",
        ctx.teamId,
        data.item.channel,
        data.item.ts,
      )[0];
      if (!entry || entry.status === "verified") return;
      const next = entry.status === "review" ? "verified" : "review";
      await status(ctx, entry.id, next);
      await ctx.slack.messages.post({
        channelId: data.item.channel,
        threadTs: data.item.ts,
        text: "Promoted " + entry.title + " to " + next,
      });
    });
  },
);
