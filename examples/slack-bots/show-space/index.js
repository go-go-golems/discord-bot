const { defineBot } = require("slack"),
  ui = require("slack/ui");
const { createStore } = require("../lib/sqlite");
const { authorized } = require("../lib/authorized");
const seeds = require("./shows.json");
const store = createStore([
  "CREATE TABLE IF NOT EXISTS shows(team TEXT,id TEXT,artist TEXT,date TEXT,doors TEXT,age TEXT,price TEXT,notes TEXT,status TEXT,channel TEXT,ts TEXT,PRIMARY KEY(team,id))",
]);
function ensure(ctx) {
  store.ensure(ctx);
  if (
    ctx.config.seedShows &&
    !store.query("SELECT id FROM shows WHERE team=? LIMIT 1", ctx.teamId).length
  ) {
    for (let i = 0; i < seeds.length; i++) {
      const s = seeds[i];
      store.exec(
        "INSERT INTO shows VALUES(?,?,?,?,?,?,?,?,?,?,?)",
        ctx.teamId,
        "seed-" + (i + 1),
        s.artist,
        s.date,
        s.doors_time,
        s.age_restriction,
        s.price,
        s.notes,
        "active",
        "",
        "",
      );
    }
  }
}
function get(ctx, id) {
  ensure(ctx);
  return store.query(
    "SELECT * FROM shows WHERE team=? AND id=?",
    ctx.teamId,
    id,
  )[0];
}
function message(show) {
  return ui
    .message(show.artist + " · " + show.date + " [" + show.status + "]")
    .block(ui.header(show.artist.slice(0, 150)))
    .block(
      ui.section(
        ui.plain(
          [
            show.date,
            "Doors: " + show.doors,
            "Age: " + show.age,
            "Price: " + show.price,
            show.notes,
            "Status: " + show.status,
          ]
            .join("\n")
            .slice(0, 3000),
        ),
      ),
    )
    .build();
}
async function canManage(ctx) {
  return authorized(
    ctx,
    ctx.config.managerIds,
    ctx.config.managerGroups,
    false,
  );
}
async function unpin(ctx, show) {
  if (show.ts) {
    try {
      await ctx.slack.pins.remove({
        channel: show.channel,
        timestamp: show.ts,
      });
    } catch (err) {
      if (err.code !== "no_pin") throw err;
    }
  }
}
async function transition(ctx, id, status) {
  const show = get(ctx, id);
  if (!show) throw Error("Show not found.");
  show.status = status;
  if (show.ts) {
    const payload = message(show);
    await ctx.slack.messages.update({
      channel: show.channel,
      ts: show.ts,
      text: payload.text,
      blocks: payload.blocks,
    });
    await unpin(ctx, show);
  }
  store.exec(
    "UPDATE shows SET status=? WHERE team=? AND id=?",
    status,
    ctx.teamId,
    id,
  );
  return show;
}
module.exports = defineBot(({ configure, command, action, view }) => {
  configure({
    name: "show-space",
    description:
      "Persistent show announcements with explicit manager authorization",
    scopes: [
      "pins:write",
      "pins:read",
      "usergroups:read",
      "users:read",
      "team:read",
    ],
    run: {
      fields: {
        dbPath: { type: "string", default: "./slack-shows.sqlite" },
        managerIds: { type: "string", default: "" },
        managerGroups: { type: "string", default: "" },
        announcementChannel: { type: "string", default: "" },
        seedShows: { type: "bool", default: true },
      },
    },
  });
  command("/upcoming", { description: "Upcoming shows" }, async (ctx) => {
    ensure(ctx);
    const rows = store.query(
      "SELECT * FROM shows WHERE team=? AND date>=? AND status='active' ORDER BY date",
      ctx.teamId,
      new Date().toISOString().slice(0, 10),
    );
    return {
      text:
        rows
          .map((s) => s.id + " · " + s.date + " · " + s.artist)
          .join("\n")
          .slice(0, 3900) || "No upcoming shows.",
    };
  });
  command(
    "/past-shows",
    { description: "Past or archived shows" },
    async (ctx) => {
      ensure(ctx);
      return {
        text:
          store
            .query(
              "SELECT * FROM shows WHERE team=? AND (date<? OR status!='active') ORDER BY date DESC",
              ctx.teamId,
              new Date().toISOString().slice(0, 10),
            )
            .map(
              (s) =>
                s.id +
                " · " +
                s.date +
                " · " +
                s.artist +
                " [" +
                s.status +
                "]",
            )
            .join("\n")
            .slice(0, 3900) || "No past shows.",
      };
    },
  );
  command("/show", { description: "Show details by ID" }, async (ctx) => {
    const show = get(ctx, ctx.text.trim());
    return show ? message(show) : { text: "Show not found." };
  });
  for (const name of ["announce", "add-show"])
    command(
      "/" + name,
      { description: "Open a show announcement form" },
      async (ctx) => {
        if (!(await canManage(ctx)))
          return {
            text: "Configure managerIds or managerGroups to authorize show management.",
          };
        let form = ui
          .modal("show.save", "Announce show")
          .metadata(ctx.channelId);
        for (const [id, label] of [
          ["artist", "Artist"],
          ["date", "Date (YYYY-MM-DD)"],
          ["doors", "Doors time"],
          ["age", "Age restriction"],
          ["price", "Price"],
          ["notes", "Notes"],
        ])
          form.input(id, label, ui.textInput("value"), {
            optional: id === "notes",
          });
        await ctx.openModal(form.submit("Publish").build());
      },
    );
  view("show.save", async (ctx) => {
    if (!(await canManage(ctx))) {
      await ctx.ack.errors({ artist: "Manager permission required." });
      return;
    }
    const show = {};
    for (const id of ["artist", "date", "doors", "age", "price", "notes"])
      show[id] = (ctx.values.text(id, "value") || "").trim();
    if (
      !/^\d{4}-\d{2}-\d{2}$/.test(show.date) ||
      !Number.isFinite(Date.parse(show.date)) ||
      new Date(show.date).toISOString().slice(0, 10) !== show.date
    ) {
      await ctx.ack.errors({ date: "Use a valid YYYY-MM-DD date." });
      return;
    }
    if (!show.artist) {
      await ctx.ack.errors({ artist: "Artist required." });
      return;
    }
    const channel = ctx.config.announcementChannel || ctx.view.privateMetadata;
    ensure(ctx);
    show.id = Date.now().toString(36);
    show.status = "active";
    await ctx.ack.accept();
    const ref = await ctx.slack.messages.post(
      Object.assign({ channelId: channel }, message(show)),
    );
    // Store the returned reference before pinning, so pin errors do not orphan the show.
    store.exec(
      "INSERT INTO shows VALUES(?,?,?,?,?,?,?,?,?,?,?)",
      ctx.teamId,
      show.id,
      show.artist,
      show.date,
      show.doors,
      show.age,
      show.price,
      show.notes,
      "active",
      ref.channelId,
      ref.ts,
    );
    await ctx.slack.pins.add({ channel: ref.channelId, timestamp: ref.ts });
    await ctx.slack.messages.ephemeral({
      channel,
      user: ctx.userId,
      text: "Posted and pinned show " + show.id,
    });
  });
  for (const [name, next] of [
    ["cancel-show", "cancelled"],
    ["archive-show", "archived"],
  ])
    command("/" + name, { description: name + " by ID" }, async (ctx) => {
      if (!(await canManage(ctx)))
        return { text: "Manager permission required." };
      try {
        return message(await transition(ctx, ctx.text.trim(), next));
      } catch (err) {
        return { text: err.message };
      }
    });
  for (const name of ["unpin-old", "archive-expired"])
    command(
      "/" + name,
      { description: "Process expired show announcements" },
      async (ctx) => {
        if (!(await canManage(ctx)))
          return { text: "Manager permission required." };
        ensure(ctx);
        const shows = store.query(
          "SELECT * FROM shows WHERE team=? AND date<? AND status='active'",
          ctx.teamId,
          new Date().toISOString().slice(0, 10),
        );
        for (const show of shows) {
          if (name === "unpin-old") await unpin(ctx, show);
          else await transition(ctx, show.id, "archived");
        }
        return { text: "Processed " + shows.length + " expired shows." };
      },
    );
  async function debug(ctx, section) {
    if (!(await canManage(ctx)))
      return { text: "Manager permission required." };
    let value = {
      teamId: ctx.teamId,
      userId: ctx.userId,
      channelId: ctx.channelId,
      manager: true,
    };
    if (section === "groups") value = await ctx.slack.usergroups.list({});
    if (section === "member")
      value = await ctx.slack.users.info({ user: ctx.userId });
    if (section === "workspace") value = await ctx.slack.workspace.info({});
    if (section === "config") value = ctx.config;
    return ui
      .message(JSON.stringify(value).slice(0, 3900))
      .block(
        ui.actions(
          "debug",
          ...[
            ["summary", "Summary"],
            ["member", "Member"],
            ["workspace", "Workspace"],
            ["config", "Config"],
            ["checks", "Checks"],
          ].map(([id, label]) => ui.button("show.debug." + id, label)),
        ),
      )
      .build();
  }
  command("/debug", { description: "Show bot diagnostics" }, (ctx) =>
    debug(ctx, "summary"),
  );
  command(
    "/debug-groups",
    { description: "Inspect workspace user groups" },
    (ctx) => debug(ctx, "groups"),
  );
  command(
    "/debug-my-permissions",
    { description: "Inspect your management permission" },
    (ctx) => debug(ctx, "member"),
  );
  for (const section of ["summary", "member", "workspace", "config", "checks"])
    action("show.debug." + section, (ctx) => debug(ctx, section));
});
