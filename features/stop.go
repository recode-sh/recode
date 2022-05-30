package features

import (
	"fmt"

	"github.com/recode-sh/recode/actions"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

type StopInput struct {
	ResolvedRepository entities.ResolvedDevEnvRepository
	PreStopHook        entities.HookRunner
}

type StopOutput struct {
	Error   error
	Content *StopOutputContent
	Stepper stepper.Stepper
}

type StopOutputContent struct {
	Cluster              *entities.Cluster
	DevEnv               *entities.DevEnv
	DevEnvAlreadyStopped bool
	SetDevEnvAsStopped   func() error
}

type StopOutputHandler interface {
	HandleOutput(StopOutput) error
}

type StopFeature struct {
	stepper             stepper.Stepper
	outputHandler       StopOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewStopFeature(
	stepper stepper.Stepper,
	outputHandler StopOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) StopFeature {

	return StopFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (s StopFeature) Execute(input StopInput) error {
	handleError := func(err error) error {
		s.outputHandler.HandleOutput(StopOutput{
			Stepper: s.stepper,
			Error:   err,
		})

		return err
	}

	devEnvName := entities.BuildDevEnvNameFromResolvedRepo(
		input.ResolvedRepository,
	)

	s.stepper.StartTemporaryStep(
		fmt.Sprintf("Stopping the development environment \"%s\"", devEnvName),
	)

	cloudService, err := s.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	recodeConfig, err := cloudService.LookupRecodeConfig(
		s.stepper,
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

	if devEnv.Status == entities.DevEnvStatusRemoving {
		return handleError(entities.ErrStopRemovingDevEnv{
			DevEnvName: devEnv.Name,
		})
	}

	if devEnv.Status == entities.DevEnvStatusCreating {
		return handleError(entities.ErrStopCreatingDevEnv{
			DevEnvName: devEnv.Name,
		})
	}

	if devEnv.Status == entities.DevEnvStatusStarting {
		return handleError(entities.ErrStopStartingDevEnv{
			DevEnvName: devEnv.Name,
		})
	}

	if input.PreStopHook != nil {
		err = input.PreStopHook.Run(
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)

		if err != nil {
			return handleError(err)
		}
	}

	devEnvAlreadyStopped := devEnv.Status == entities.DevEnvStatusStopped

	if !devEnvAlreadyStopped {

		err = actions.StopDevEnv(
			s.stepper,
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)

		if err != nil {
			return handleError(err)
		}
	}

	setDevEnvAsStopped := func() error {
		devEnv.Status = entities.DevEnvStatusStopped

		return actions.UpdateDevEnvInConfig(
			s.stepper,
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)
	}

	return s.outputHandler.HandleOutput(StopOutput{
		Stepper: s.stepper,
		Content: &StopOutputContent{
			Cluster:              cluster,
			DevEnv:               devEnv,
			DevEnvAlreadyStopped: devEnvAlreadyStopped,
			SetDevEnvAsStopped:   setDevEnvAsStopped,
		},
	})
}
