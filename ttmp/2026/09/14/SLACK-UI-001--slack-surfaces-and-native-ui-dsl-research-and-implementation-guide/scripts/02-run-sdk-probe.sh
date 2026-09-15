#!/bin/sh
set -eu
# Run from discord-bot root; uses only the existing top-level module.
probe_dir=$(mktemp -d /tmp/slack-ui-sdk.XXXXXX)
trap 'rm -rf "$probe_dir"' EXIT
cp ttmp/2026/09/14/SLACK-UI-001--slack-surfaces-and-native-ui-dsl-research-and-implementation-guide/scripts/02-sdk-block-probe.go.txt "$probe_dir/main.go"
sh ttmp/2026/09/10/DISCORD-SLACK-001--add-slack-support-to-discord-bot/scripts/04-go-offline.sh run "$probe_dir/main.go"
