# Cloud Agent

[![license](https://img.shields.io/badge/License-MIT-brightgreen.svg?style=for-the-badge)](https://github.com/pPrecel/cloudagent/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pPrecel/cloudagent?style=for-the-badge)](https://goreportcard.com/report/github.com/pPrecel/cloudagent)
[![build](https://img.shields.io/github/workflow/status/pPrecel/cloudagent/build?style=for-the-badge)](https://github.com/pPrecel/cloudagent/actions/workflows/build.yml)
[![Coverage](https://img.shields.io/coveralls/github/pPrecel/cloudagent?style=for-the-badge)](https://coveralls.io/github/pPrecel/cloudagent)

The simple and easy-to-use program is designed to watch user activity and possible orphan clusters for Cloud Providers:

- Gardener
- GCP (work in progress)

This application is created with a view to using it as [the tmux](https://github.com/tmux/tmux) status. To fulfill this criterion the procedure of getting resources from the gardener is separated and is in the second service which serves the UNIX socket that is used by the first one. This architecture allows not to block the main tmux process during calling the right gardener endpoint.

## Installation

### Github Release

Visit the [releases page](https://github.com/pPrecel/cloudagent/releases) to download one of the pre-built binaries for your platform.

### Homebrew

1. Use Homebrew to install `cloudagent`:

    ```bash
    brew install pPrecel/tap/cloudagent
    ```

    or

    ```bash
    brew tap pPrecel/tap
    brew install cloudagent
    ```

2. Extend configuration ([see also](./docs/configuration-file.md)) in the `~/.cloudagent.conf.yaml` location.

3. Start the `cloudagent` service:

    ```bash
    brew services start cloudagent
    ```

## Make use of it

After installing, configuring, and starting the service process, cloudagent would fetch info from given clouds after ~60sec from start. To check that program is working correctly use:

```bash
cloudagent state
```

> **NOTE:** I really recommend you to use `--help` for all commands to familiarize with all functionalities. Logs from running service are stored in `/tmp/cloudagent/cloudagent.stdout` and you can use for example `cat` program to read all of them.

## Integration with tmux

To add this application to tmux put the line below in the `~/.tmux.conf` file:

```text
set -ag status-right ' #(cloudagent state --createdBy <OWNER_NAME> -o text) '
```

## Contribution

Any help and contribution would be welcome. If you want to help with creating feature requests, logging bugs, or/and working on any existing ticket, go to the [issues](https://github.com/pPrecel/cloudagent/issues) tab. If you decided to look into code, [this document](./docs/local-development.md) would be helpful.
