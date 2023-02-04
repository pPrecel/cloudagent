# Cloud Agent

[![license](https://img.shields.io/badge/License-MIT-brightgreen.svg?style=for-the-badge)](https://github.com/pPrecel/cloudagent/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pPrecel/cloudagent?style=for-the-badge)](https://goreportcard.com/report/github.com/pPrecel/cloudagent)
[![build](https://img.shields.io/github/actions/workflow/status/pPrecel/cloudagent/tests-build.yml?style=for-the-badge)](https://github.com/pPrecel/cloudagent/actions/workflows/build.yml)
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

2. Start the `cloudagent` service:

    ```bash
    brew services start cloudagent
    ```

### Manual

1. Go to the [Releases](https://github.com/pPrecel/cloudagent/releases/latest), and download the right cloudagent version for your system.

> **NOTE:** We don't fully support any Linux package managers ( except brew ). That means you must on your own craft a service background configuration ( Daemon ).

## Make use of it

After installing and starting the service process you can pass configuration to immediately say agent which clouds should he observe (read more [here](./docs/configuration-file.md)). To check that program is working correctly use:

```bash
cloudagent state
```

> **NOTE:** I really recommend you to use `--help` for all commands to familiarize with all functionalities. Logs from running service are stored in `/tmp/cloudagent/cloudagent.stdout` and you can use for example `cat` program to read all of them.

## Integrations

The cloudagent may be integrated with other tools like the tmux to observe clusters or the vscode to easily extend the `.cloudagent.conf.yml` file ( read more [here](./docs/integrations.md) ).

## Contribution

Any help and contribution would be welcome. If you want to help with creating feature requests, logging bugs, or/and working on any existing ticket, go to the [issues](https://github.com/pPrecel/cloudagent/issues) tab. If you decided to look into code, [this document](./docs/local-development.md) would be helpful.
