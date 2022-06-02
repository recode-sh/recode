package features

import (
	"fmt"

	"github.com/recode-sh/recode/actions"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

type RemoveInput struct {
	ResolvedRepository entities.ResolvedDevEnvRepository
	PreRemoveHook      entities.HookRunner
	ForceRemove        bool
	ConfirmRemove      func() (bool, error)
}

type RemoveOutput struct {
	Error   error
	Content *RemoveOutputContent
	Stepper stepper.Stepper
}

type RemoveOutputContent struct {
	Cluster *entities.Cluster
	DevEnv  *entities.DevEnv
}

type RemoveOutputHandler interface {
	HandleOutput(RemoveOutput) error
}

type RemoveFeature struct {
	stepper             stepper.Stepper
	outputHandler       RemoveOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewRemoveFeature(
	stepper stepper.Stepper,
	outputHandler RemoveOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) RemoveFeature {

	return RemoveFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (r RemoveFeature) Execute(input RemoveInput) error {
	handleError := func(err error) error {
		r.outputHandler.HandleOutput(RemoveOutput{
			Stepper: r.stepper,
			Error:   err,
		})

		return err
	}

	devEnvName := entities.BuildDevEnvNameFromResolvedRepo(
		input.ResolvedRepository,
	)

	step := fmt.Sprintf("Removing the development environment \"%s\"", devEnvName)
	r.stepper.StartTemporaryStep(step)

	cloudService, err := r.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	recodeConfig, err := cloudService.LookupRecodeConfig(
		r.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := recodeConfig.GetCluster(clusterName)

	if err != nil {
		return handleError(err)
	}

	devEnv, err := recodeConfig.GetDevEnv(cluster.Name, devEnvName)

	if err != nil {
		return handleError(err)
	}

	if !input.ForceRemove && input.ConfirmRemove != nil {
		r.stepper.StopCurrentStep()

		confirmed, err := input.ConfirmRemove()

		if err != nil {
			return handleError(err)
		}

		if !confirmed {
			return nil
		}

		r.stepper.StartTemporaryStep(step)
	}

	err = actions.RemoveDevEnv(
		r.stepper,
		cloudService,
		recodeConfig,
		cluster,
		devEnv,
		input.PreRemoveHook,
	)

	if err != nil {
		return handleError(err)
	}

	return r.outputHandler.HandleOutput(RemoveOutput{
		Stepper: r.stepper,
		Content: &RemoveOutputContent{
			Cluster: cluster,
			DevEnv:  devEnv,
		},
	})
}
