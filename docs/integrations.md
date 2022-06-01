# Integrations

This article is aimed to describe the most valuable integrations with external tools. This does not mean that other integrations are impossible and do not exist. The only blocker is your imagination.

## Integration with tmux

To add this application to tmux put the line below in the `~/.tmux.conf` file:

```text
set -ag status-right ' #(cloudagent state --createdBy <OWNER_NAME> -o text) '
```

## Integration with vscode

To support the `.cloudagent.conf.yml` schema file by default generate schema and then add schema file path to the vscode settings:

1. Generate schema:

    ```bash
    cloudagent config schema > <HOME_PATH>/.cloudagent.conf.schema
    ```

2. Extend the vscode settings file by adding lines below:

    ```json
    "yaml.schemas": {
        "<HOME_PATH>/.cloudagent.conf.schema": "*.cloudagent.conf.yaml"
    }
    ```
