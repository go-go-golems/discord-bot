Version `4.1.0` of the Slack CLI has landed!

- We've added a `slack api <method>` command, which allows you to call any API method directly from the Slack CLI. For example:
	- `slack api chat.postMessage channel=C123 text="hi"`
	- `slack api chat.postMessage '{"channel":"C123","text":"hi"}'`
- Interactive forms are now screen reader friendly when the `ACCESSIBLE` environment variable is set, or when using the `--accessible` flag.
- We've added a new `env init` command to copy a template placeholder `.env.sample` or `.env.example` file to the `.env` file of your project.
- We've added support for the `settings.is_mcp_enabled: <bool>` manifest property to enable or disable the Slack MCP Server. This causes cached values to expire and the next run attempt may surface a warning noting that app settings have changed since the last update.
- We've added support for optional OAuth scopes (`bot_optional` and `user_optional`) in the app manifest schema.
- We've added a `SLACK_CLI_APP_ICON_PATH` environment variable to override the icon file path used during app install and local run. When set, the specified path takes priority over the manifest `icon` field and the default icon file search. A warning is displayed if the configured path does not exist.
- Icon auto-detection now supports `.png`, `.jpg`, `.jpeg`, and `.gif` formats in `assets/` and project root. The `assets/` directory is checked before project root; `.png` is preferred when multiple formats exist.
- App manifest display names now preserve the original app name instead of using the kebab-case directory name during `slack create`.
- We've added verbose debug logging and retry support to icon upload API calls.
- We capitalized "App" in the `slack create` category selection menu titles for more consistent title casing.
- We've added more polish to outputs from the `slack create` command section headers.
- The `slack create` command preserves path separators when a nested path is provided as the project directory. The `--name` flag now only overrides the app's display name in the manifest without replacing the path argument. For example:
	```sh
	$ slack create path/to/my-app --name "My App"
	# Creates directory at ./path/to/my-app/ with manifest name "My App"
	```
- We fixed the S3 `Content-MD5` header when deploying a Deno app to Slack. We now remove the MD5 hash computation and `Content-MD5` header from S3 presigned POST uploads to avoid an error.
- The Slack CLI installer now outputs the pinned version provided if a particular version is provided, rather than a version and blank space.
- The Slack CLI now avoids uploading icons for non-hosted apps that do not have the `set-icon` experiment enabled.