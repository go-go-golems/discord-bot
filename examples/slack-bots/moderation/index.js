const { defineBot } = require("slack");
const { authorized } = require("../lib/authorized");
const pretty = (value) => ({
  text: JSON.stringify(value, null, 2).slice(0, 3900),
});
module.exports = defineBot(({ configure, command, event }) => {
  configure({
    name: "moderation",
    description: "Authorized Slack message/channel/user administration",
    scopes: [
      "channels:read",
      "channels:history",
      "channels:manage",
      "channels:write.topic",
      "users:read",
      "usergroups:read",
      "usergroups:write",
      "pins:read",
      "pins:write",
      "team:read",
    ],
    run: {
      fields: {
        moderatorIds: { type: "string", default: "" },
        moderatorGroups: { type: "string", default: "" },
        enableWorkspaceRemoval: { type: "bool", default: false },
        useUserToken: { type: "bool", default: false },
      },
    },
  });
  const guarded = (handler) => async (ctx) => {
    if (
      !(await authorized(
        ctx,
        ctx.config.moderatorIds,
        ctx.config.moderatorGroups,
        false,
      ))
    )
      return {
        text: "Moderator permission required. Configure moderatorIds or moderatorGroups.",
      };
    try {
      return await handler(ctx, ctx.text.trim().split(/\s+/).filter(Boolean));
    } catch (err) {
      return { text: err.code ? err.operation + ": " + err.code : err.message };
    }
  };
  const add = (name, description, handler) =>
    command("/mod-" + name, { description }, guarded(handler));
  add("summary", "Moderation capabilities", async () => ({
    text: "Messages, pins, topics, members and user groups. Workspace removal requires an explicit setting and Enterprise user token. Slack has no direct Discord timeout, ban/unban or slowmode counterpart.",
  }));
  add("guidelines", "Moderation guidelines", async () => ({
    text: "Inspect targets before mutations. Bot tokens delete only messages they authored. User groups are membership lists, not permission roles. Removing a workspace member is distinct from removing a channel member.",
  }));
  add("list-messages", "Recent messages [limit] [cursor]", async (ctx, args) =>
    pretty(
      await ctx.slack.conversations.history({
        channel: ctx.channelId,
        limit: Math.max(1, Math.min(100, Number(args[0] || 20))),
        cursor: args[1] || "",
      }),
    ),
  );
  add("fetch-message", "Fetch a message by timestamp", async (ctx, args) => {
    if (!args[0]) return { text: "Message timestamp required." };
    const result = await ctx.slack.conversations.history({
      channel: ctx.channelId,
      latest: args[0],
      inclusive: true,
      limit: 1,
    });
    const message = (result.messages || []).find((m) => m.ts === args[0]);
    return message ? pretty(message) : { text: "Message not found." };
  });
  add("pin", "Pin message timestamp", async (ctx, args) =>
    pretty(
      await ctx.slack.pins.add({ channel: ctx.channelId, timestamp: args[0] }),
    ),
  );
  add("unpin", "Unpin message timestamp", async (ctx, args) =>
    pretty(
      await ctx.slack.pins.remove({
        channel: ctx.channelId,
        timestamp: args[0],
      }),
    ),
  );
  add("list-pins", "List channel pins", async (ctx) =>
    pretty(await ctx.slack.pins.list({ channel: ctx.channelId })),
  );
  add(
    "bulk-delete",
    "Delete bot-authored messages by timestamps",
    async (ctx, args) => {
      if (!args.length || args.length > 100)
        return { text: "Supply 1–100 message timestamps." };
      let deleted = 0;
      for (const ts of args) {
        try {
          await (
            ctx.config.useUserToken
              ? ctx.slack.messages.deleteAsUser
              : ctx.slack.messages.delete
          )({ channel: ctx.channelId, ts });
          deleted++;
        } catch (err) {
          return {
            text:
              "Deleted " +
              deleted +
              " messages; stopped at " +
              ts +
              ": " +
              (err.code || err.message),
          };
        }
      }
      return { text: "Deleted " + deleted + " messages." };
    },
  );
  add("fetch-channel", "Fetch current channel", async (ctx) =>
    pretty(await ctx.slack.conversations.info({ channel: ctx.channelId })),
  );
  add("set-topic", "Set current channel topic", async (ctx) =>
    pretty(
      await ctx.slack.conversations.setTopic({
        channel: ctx.channelId,
        topic: ctx.text,
      }),
    ),
  );
  add("fetch-workspace", "Fetch workspace details", async (ctx) =>
    pretty(await ctx.slack.workspace.info({})),
  );
  // Register native group names: Slack user groups do not confer Discord role permissions.
  add("list-groups", "List workspace user groups", async (ctx) =>
    pretty(await ctx.slack.usergroups.list({ include_users: true })),
  );
  add("fetch-group", "Fetch user-group ID", async (ctx, args) => {
    const result = await ctx.slack.usergroups.list({ include_users: true });
    return pretty(
      (result.usergroups || []).find((group) => group.id === args[0]) || {
        error: "Group not found",
      },
    );
  });
  add("fetch-member", "Fetch Slack user ID", async (ctx, args) =>
    pretty(await ctx.slack.users.info({ user: args[0] })),
  );
  add("list-members", "List workspace members [cursor]", async (ctx, args) =>
    pretty(await ctx.slack.users.list({ limit: 100, cursor: args[0] || "" })),
  );
  add("add-group-member", "group-id user-id", async (ctx, args) => {
    if (args.length !== 2)
      return { text: "Usage: /mod-add-group-member group-id user-id" };
    const result = await ctx.slack.usergroups.members({ usergroup: args[0] });
    const users = Array.from(new Set((result.users || []).concat([args[1]])));
    await (
      ctx.config.useUserToken
        ? ctx.slack.usergroups.setMembersAsUser
        : ctx.slack.usergroups.setMembers
    )({
      usergroup: args[0],
      users: users.join(","),
    });
    return {
      text: "Added group member; preserved " + users.length + " total members.",
    };
  });
  add(
    "kick-channel",
    "Remove user ID from current channel",
    async (ctx, args) =>
      pretty(
        await ctx.slack.conversations.kick({
          channel: ctx.channelId,
          user: args[0],
        }),
      ),
  );
  add(
    "remove-workspace-user",
    "Enterprise: remove user ID from this workspace",
    async (ctx, args) => {
      if (!ctx.config.enableWorkspaceRemoval)
        return {
          text: "Workspace removal is disabled. Enable explicitly and import an Enterprise user token.",
        };
      if (!args[0]) return { text: "User ID required." };
      return pretty(
        await ctx.slack.admin.removeUser({
          team_id: ctx.teamId,
          user_id: args[0],
        }),
      );
    },
  );
  for (const name of [
    "message",
    "message_changed",
    "message_deleted",
    "user_change",
    "reaction_added",
    "reaction_removed",
    "member_joined_channel",
    "member_left_channel",
    "team_join",
  ])
    event(name, async (ctx) => {
      ctx.log.info(
        "Moderation event " +
          name +
          " user=" +
          ctx.userId +
          " channel=" +
          ctx.channelId,
      );
    });
});
