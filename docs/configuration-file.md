# Configuration file

The Cloud Agent can be fully configured by creating the `.cloudagent.conf.yaml` file in the `${HOME}` directory. The file should contain information about the clouds you want to observe.

The Cloud Agent has functionality of auto-detecting any change around the configuration file. It means that after every change in the configuration file, program will automatically fetch all changes.

> **NOTE:** The configuration file can be easily maintained using JSON schema ( read more [here](./integrations.md) ).

The sample config file:

```yaml
persistentSpec: "@every 2s"
gardenerProjects:
  - namespace: "n1"
    kubeconfigPath: "/path1"
  - namespace: "n2"
    kubeconfigPath: "/path1"
  - namespace: "n3"
    kubeconfigPath: "/path2"
gcpProjects: []
```

## Parameters

See all parameter descriptions:

| Field  | Description |
|-|-|
| **persistentSpec** | A cron extension represents a set of times. More info [here](https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format). Also possible is to set this field to `on-demand` to run the agent in a mode without caching. |
| **gardenerProjects** | List of gardener projects to observe. |
| **gardenerProjects[].namespace** | Project namespace. |
| **gardenerProjects[].kubeconfigPath** | Path to projects kubeconfig. |
| **gcpProjects** | List of GCP projects to observe. YET NOT IMPLEMENTED. |
