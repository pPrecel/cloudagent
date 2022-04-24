# Cloud Agent

[![license](https://img.shields.io/badge/License-MIT-brightgreen.svg?style=for-the-badge)](https://github.com/pPrecel/cloud-agent/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pPrecel/cloud-agent?style=for-the-badge)](https://goreportcard.com/report/github.com/pPrecel/cloud-agent)
[![build](https://img.shields.io/github/workflow/status/pPrecel/cloud-agent/build?style=for-the-badge)](https://github.com/pPrecel/cloud-agent/actions/workflows/build.yml)

The simple and easy-to-use program designed to watch user activity and possible orphan clusters for Cloud Providers:

- Gardener
- GCP (work in progress)

This application is created with a view to using it as [the tmux](https://github.com/tmux/tmux) status. To fulfill this criterion the procedure of getting resources from the gardener is separated and is in the second service which serves the UNIX socket that is used by the first one. This architecture allows not to block the main tmux process during calling the right gardener endpoint.

## Prerequisites

## Installation

1. Build program:

    ```bash
    make build
    ```

2. Add program to PATH and install it as a system agent:

    ```bash
    make ln-to-path
    make install-agent kubeconfigPath=<KUBECONFIG_PATH> namespace=<NAMESPACE>
    ```

    > **NOTE:** for local development or need to get more informations from the agent you can pass more arguments to the `make install-agent` command like: `other_flags=--agentVerbose`.

3. Check if program works by getting its logs:

    ```bash
    tail /tmp/cloud-agent.stdout
    ```

4. After waiting ~60 seconds for first iteration of the watcher you can get cluster state:

    ```bash
    cloudagent state --createdBy <OWNER_NAME>
    ```

## Un-installation

1. Remove application from path and remove system agent:

    ```bash
    make rm-from-path
    make uninstall-agent
    ```

## Integration with tmux

To add this application to tmux put line below in the `~/.tmux.conf` file:

```text
set -ag status-right ' #(cloudagent state --createdBy <OWNER_NAME> -o text) '
```
