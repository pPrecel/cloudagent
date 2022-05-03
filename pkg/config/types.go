package config

type Config struct {
	PersistentSpec   string            `yaml:"persistentSpec"`
	GardenerProjects []GardenerProject `yaml:"gardenerProjects"`
	GCPProjects      []GCPProject      `yaml:"gcpProjects"`
}

type GardenerProject struct {
	Namespace      string `yaml:"namespace"`
	KubeconfigPath string `yaml:"kubeconfigPath"`
}

type GCPProject struct {
	// TODO: not implemented
}
