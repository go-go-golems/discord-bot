package slackbot

import "context"

// OperationService is a bounded set of Slack Web API capabilities. Parameters
// and results use Slack's documented wire keys, detached from SDK objects.
type OperationService interface {
	Call(context.Context, string, map[string]any) (map[string]any, error)
}

// OperationMethods returns a new map so callers cannot mutate the allowlist.
func OperationMethods() map[string]string {
	return map[string]string{
		"admin.removeUser":            "admin.users.remove",
		"messages.deleteAsUser":       "chat.delete",
		"usergroups.setMembersAsUser": "usergroups.users.update",
		"files.upload":                "files.uploadExternal",
		"messages.update":             "chat.update", "messages.delete": "chat.delete", "messages.ephemeral": "chat.postEphemeral", "messages.permalink": "chat.getPermalink",
		"conversations.history": "conversations.history", "conversations.replies": "conversations.replies", "conversations.info": "conversations.info", "conversations.list": "conversations.list", "conversations.members": "conversations.members", "conversations.join": "conversations.join", "conversations.leave": "conversations.leave", "conversations.setTopic": "conversations.setTopic", "conversations.kick": "conversations.kick", "conversations.archive": "conversations.archive",
		"users.info": "users.info", "users.list": "users.list",
		"usergroups.list": "usergroups.list", "usergroups.members": "usergroups.users.list", "usergroups.setMembers": "usergroups.users.update",
		"pins.add": "pins.add", "pins.remove": "pins.remove", "pins.list": "pins.list",
		"reactions.add": "reactions.add", "reactions.remove": "reactions.remove", "reactions.get": "reactions.get",
		"workspace.info": "team.info",
	}
}
