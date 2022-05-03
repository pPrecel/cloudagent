# Cloud Agent

[![license](https://img.shields.io/badge/License-MIT-brightgreen.svg?style=for-the-badge)](https://github.com/pPrecel/cloud-agent/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pPrecel/cloud-agent?style=for-the-badge)](https://goreportcard.com/report/github.com/pPrecel/cloud-agent)
[![build](https://img.shields.io/github/workflow/status/pPrecel/cloud-agent/build?style=for-the-badge)](https://github.com/pPrecel/cloud-agent/actions/workflows/build.yml)
[![Coverage](https://img.shields.io/coveralls/github/pPrecel/cloud-agent?style=for-the-badge)](https://coveralls.io/github/pPrecel/cloud-agent)

The simple and easy-to-use program is designed to watch user activity and possible orphan clusters for Cloud Providers:

- Gardener
- GCP (work in progress)

This application is created with a view to using it as [the tmux](https://github.com/tmux/tmux) status. To fulfill this criterion the procedure of getting resources from the gardener is separated and is in the second service which serves the UNIX socket that is used by the first one. This architecture allows not to block the main tmux process during calling the right gardener endpoint.

## Installation

1. Verify and build program:

    ```bash
    make verify
    make build
    ```

2. Create configuration file ([see also](./docs/configuration-file.md)) in the `${HOME}/.cloudagent.conf.yaml` location.

3. Add the program to PATH and install it as a system agent:

    ```bash
    make ln-to-path
    make install-agent
    ```

    > **NOTE:** for local development or need to get more information from the agent you can pass more arguments to the `make install-agent` command like: `other_flags=--agentVerbose`.

4. Check if the program works by getting its logs:

    ```bash
    tail /tmp/cloud-agent.stdout
    ```

5. After waiting ~60 seconds for the first iteration of the watcher you can get the cluster state:

    ```bash
    cloudagent state
    ```

## Un-installation

1. Remove the application from the path and remove system agent:

    ```bash
    make rm-from-path
    make uninstall-agent
    ```

## Integration with tmux

To add this application to tmux put line below in the `~/.tmux.conf` file:

```text
set -ag status-right ' #(cloudagent state --createdBy <OWNER_NAME> -o text) '
```
