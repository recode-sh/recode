package entities

type HookRunner interface {
	Run(
		cloudService CloudService,
		config *Config,
		cluster *Cluster,
		devEnv *DevEnv,
	) error
}
