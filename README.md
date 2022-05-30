# Recode

This repository contains all the entities, actions and features used by all the other packages (including the [Recode CLI](https://github.com/recode-sh/cli)). 

*(For the clean architecture aficionados, we are in the innermost circle (ie: the entities one).)*

## Table of contents
- [Requirements](#requirements)
- [Usage](#usage)
- [The future](#the-future)
- [License](#license)

## Requirements

This repository makes use of go generics and, as a result, needs to have `go >= 1.18` as a requirement.

## Usage

**This repository is not meant to be used standalone (you could see that there is no `main.go` file). It is only meant to be imported by other packages**.

As an example, all the cloud providers added to the [Recode CLI](https://github.com/recode-sh/cli) need to conform to the `CloudService` interface:

```go
// entities/cloud_service.go
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
```

## The future

This project is **100% community-driven**, meaning that except for bug fixes <ins>**no more features will be added**</ins>. 

The only features that will be added are the ones that will be [posted as an issue](https://github.com/recode-sh/cli/issues/new) and that will receive a significant amount of upvotes **(>= 10 currently)**.

## License

Recode is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
