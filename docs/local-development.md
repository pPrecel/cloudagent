# Local development

This repository is based of Makefile which contains most important targets to work with this repository. To get more info read sections below or simply look info [Makefile](../Makefile).

## Validate and build

To run all checks, linters and build use:

```bash
make verify
make build
```

## Manual installation

Installation of the cloudagent during development is based on starting server manually. If you need to use cloudagent in more gentle way, try to link it to the PATH:

```bash
make ln-to-path
```

> **NOTE:** To remove it from path you can use `rm-from-path` target.

To start server create configuration file first ([read more here](./configuration-file.md)) and then run process in your terminal:

```bash
cloudagent-dev serve
```

## Re-generate protobuf

Communication between cloudagent server and client is being implemented by using the `unix-socket` and `gRPC` protocol. If your change touches the `proto` folder then you have to regenerate protobuf package. To do it use:

```bash
make protobuf
```
