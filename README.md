# Gardener Agent

The simple and easy-to-use program designed to watch user activity and possible orphan resources in [the gardener](https://github.com/gardener/gardener) namespace.

This application is created with a view to using it as [the tmux](https://github.com/tmux/tmux) status. To fulfill this criterion the procedure of getting resources from the gardener is separated and is in the second service which serves the UNIX socket that is used by the first one. This architecture allows not to block the main tmux process during calling the right gardener endpoint.

## Prerequisites

## Installation

1. Build program:

    ```bash
    make build
    ```

2. Add application to PATH:

    ```bash
    ln -s $(pwd)/.out/gardenagent /usr/local/bin/gardenagent
    ```

    > **NOTE:** you can copy binary instead of linking it by replacing the `ln -s` command with the`cp`.
