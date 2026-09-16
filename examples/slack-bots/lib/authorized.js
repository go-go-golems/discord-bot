function ids(raw) {
  return String(raw || "")
    .split(/[\s,;]+/)
    .filter(Boolean);
}
async function authorized(ctx, users, groups, allowUnconfigured) {
  const userIds = ids(users),
    groupIds = ids(groups);
  if (!userIds.length && !groupIds.length) return !!allowUnconfigured;
  if (userIds.includes(ctx.userId)) return true;
  for (const usergroup of groupIds) {
    const result = await ctx.slack.usergroups.members({ usergroup });
    if ((result.users || []).includes(ctx.userId)) return true;
  }
  return false;
}
module.exports = { ids, authorized };
