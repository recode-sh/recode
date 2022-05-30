package entities

import "github.com/recode-sh/recode/stepper"

type CloudService interface {
	CreateRecodeConfigStorage(stepper.Stepper) error
	RemoveRecodeConfigStorage(stepper.Stepper) error

	LookupRecodeConfig(stepper.Stepper) (*Config, error)
	SaveRecodeConfig(stepper.Stepper, *Config) error

	CreateCluster(stepper.Stepper, *Config, *Cluster) error
	RemoveCluster(stepper.Stepper, *Config, *Cluster) error

	CheckInstanceTypeValidity(stepper.Stepper, string) error

	CreateDevEnv(stepper.Stepper, *Config, *Cluster, *DevEnv) error
	RemoveDevEnv(stepper.Stepper, *Config, *Cluster, *DevEnv) error

	StartDevEnv(stepper.Stepper, *Config, *Cluster, *DevEnv) error
	StopDevEnv(stepper.Stepper, *Config, *Cluster, *DevEnv) error
}

type CloudServiceBuilder interface {
	Build() (CloudService, error)
}
