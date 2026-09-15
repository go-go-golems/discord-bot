---
Title: "06 surfaces split view"
Ticket: SLACK-UI-001
Status: active
Topics: [slack, architecture]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: ["https://docs.slack.dev/surfaces/split-view/"]
Summary: "Official Slack documentation source capture; see sources catalog for extraction details."
LastUpdated: 2026-09-14T23:00:00-04:00
WhatFor: "Research evidence."
WhenToUse: "Check the original API contract."
---

The split view surface allows users to initiate private conversations with agents side by side with the rest of Slack, in the flow of work. Use this view in conjunction with other app features conducive to [AI integration](https://docs.slack.dev/ai/developing-agents); for example, suggested prompts, loading states, and app threads.

[Create an app](https://api.slack.com/apps?new_app=1)

![Image of split view](https://docs.slack.dev/img/guides/ai_container/splitview.png)

## Enabling split view

Once you've created an app in the [app settings](https://api.slack.com/apps), navigate to **Agents** in the side bar, then click the toggle to turn the feature on. Turning on this feature also enables a top bar entry point.

A prior implementation of this feature offered a slightly different messaging experience. See [Migrating to the Agent messaging experience](https://docs.slack.dev/ai/migrating-to-agent-messaging) to update your app to use the new agent messaging experience.

Read more about other AI features in apps in the guide to [Developing agents](https://docs.slack.dev/ai/developing-agents).