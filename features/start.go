package features

import (
	"errors"
	"fmt"

	"github.com/recode-sh/recode/actions"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

type StartInput struct {
	InstanceType             string
	DevEnvRebuildAsked       bool
	ResolvedDevEnvUserConfig entities.ResolvedDevEnvUserConfig
	ResolvedRepository       entities.ResolvedDevEnvRepository
	ForceDevEnvRevuild       bool
	ConfirmDevEnvRebuild     func() (bool, error)
}

type StartOutput struct {
	Error   error
	Content *StartOutputContent
	Stepper stepper.Stepper
}

type StartOutputContent struct {
	CloudService       entities.CloudService
	RecodeConfig       *entities.Config
	Cluster            *entities.Cluster
	DevEnv             *entities.DevEnv
	DevEnvCreated      bool
	DevEnvStarted      bool
	SetDevEnvAsStarted func() error
	DevEnvRebuildAsked bool
}

type StartOutputHandler interface {
	HandleOutput(StartOutput) error
}

type StartFeature struct {
	stepper             stepper.Stepper
	outputHandler       StartOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewStartFeature(
	stepper stepper.Stepper,
	outputHandler StartOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) StartFeature {

	return StartFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (s StartFeature) Execute(input StartInput) error {
	handleError := func(err error) error {
		s.outputHandler.HandleOutput(StartOutput{
			Stepper: s.stepper,
			Error:   err,
		})

		return err
	}

	devEnvName := entities.BuildDevEnvNameFromResolvedRepo(
		input.ResolvedRepository,
	)

	step := fmt.Sprintf("Starting the development environment \"%s\"", devEnvName)

	if input.DevEnvRebuildAsked {
		step = fmt.Sprintf("Rebuilding the development environment \"%s\"", devEnvName)
	}

	s.stepper.StartTemporaryStep(step)

	cloudService, err := s.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	if !input.DevEnvRebuildAsked {
		err = cloudService.CheckInstanceTypeValidity(
			s.stepper,
			input.InstanceType,
		)

		if err != nil {
			return handleError(err)
		}
	}

	recodeConfig, err := cloudService.LookupRecodeConfig(
		s.stepper,
	)

	if err != nil &&
		(!errors.Is(err, entities.ErrRecodeNotInstalled) || input.DevEnvRebuildAsked) {

		return handleError(err)
	}

	if recodeConfig == nil { // Recode not installed

		s.stepper.StartTemporaryStep("Installing Recode")

		recodeConfig = entities.NewConfig()

		err = actions.InstallRecode(
			s.stepper,
			cloudService,
			recodeConfig,
		)

		if err != nil {
			return handleError(err)
		}
	}

	clusterName := entities.DefaultClusterName
	cluster, err := recodeConfig.GetCluster(clusterName)

	if err != nil &&
		(!errors.As(err, &entities.ErrClusterNotExists{}) || input.DevEnvRebuildAsked) {

		return handleError(err)
	}

	if cluster == nil || cluster.Status == entities.ClusterStatusCreating {

		/* Cluster not exists or still
		in creating state after error */

		s.stepper.StartTemporaryStep("Creating default cluster")

		if cluster == nil {
			// Multiple clusters are not implemented for now
			isDefaultCluster := true

			cluster = entities.NewCluster(
				clusterName,
				input.InstanceType,
				isDefaultCluster,
			)
		}

		err = actions.CreateCluser(
			s.stepper,
			cloudService,
			recodeConfig,
			cluster,
		)

		if err != nil {
			return handleError(err)
		}
	}

	devEnv, err := recodeConfig.GetDevEnv(
		cluster.Name,
		devEnvName,
	)

	if err != nil &&
		(!errors.As(err, &entities.ErrDevEnvNotExists{}) || input.DevEnvRebuildAsked) {

		return handleError(err)
	}

	if devEnv != nil && devEnv.Status == entities.DevEnvStatusRemoving {
		return handleError(entities.ErrStartRemovingDevEnv{
			DevEnvName: devEnv.Name,
		})
	}

	if devEnv != nil && devEnv.Status == entities.DevEnvStatusStopping {
		return handleError(entities.ErrStartStoppingDevEnv{
			DevEnvName: devEnv.Name,
		})
	}

	if input.DevEnvRebuildAsked &&
		input.ConfirmDevEnvRebuild != nil &&
		!input.ForceDevEnvRevuild {

		s.stepper.StopCurrentStep()

		confirmed, err := input.ConfirmDevEnvRebuild()

		if err != nil {
			return handleError(err)
		}

		if !confirmed {
			return nil
		}

		s.stepper.StartTemporaryStep(step)
	}

	devEnvCreated := false

	if devEnv == nil || devEnv.Status == entities.DevEnvStatusCreating {

		/* Dev env not exists or still
		in creating state after error */

		if devEnv == nil {
			devEnv = entities.NewDevEnv(
				devEnvName,
				input.InstanceType,
				input.ResolvedDevEnvUserConfig,
				input.ResolvedRepository,
			)
		}

		err = actions.CreateDevEnv(
			s.stepper,
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)

		if err != nil {
			return handleError(err)
		}

		devEnvCreated = true
	}

	devEnvStarted := false

	if devEnv.Status == entities.DevEnvStatusStopped {
		err = actions.StartDevEnv(
			s.stepper,
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)

		if err != nil {
			return handleError(err)
		}

		devEnvStarted = true
	}

	// Current step is the last ended infrastructure step.
	// Better UX if we reset to main step here given that
	// the next steps (in GRPC agent) may take some time to start.
	s.stepper.StartTemporaryStep(step)

	setDevEnvAsStarted := func() error {
		devEnv.Status = entities.DevEnvStatusStarted

		return actions.UpdateDevEnvInConfig(
			s.stepper,
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)
	}

	return s.outputHandler.HandleOutput(StartOutput{
		Stepper: s.stepper,
		Content: &StartOutputContent{
			CloudService:       cloudService,
			RecodeConfig:       recodeConfig,
			Cluster:            cluster,
			DevEnv:             devEnv,
			DevEnvCreated:      devEnvCreated,
			DevEnvStarted:      devEnvStarted,
			SetDevEnvAsStarted: setDevEnvAsStarted,
			DevEnvRebuildAsked: input.DevEnvRebuildAsked,
		},
	})
}
