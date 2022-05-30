package features

import (
	"errors"

	"github.com/recode-sh/recode/actions"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

type UninstallInput struct {
	SuccessMessage            string
	AlreadyUninstalledMessage string
}

type UninstallOutput struct {
	Error   error
	Content *UninstallOutputContent
	Stepper stepper.Stepper
}

type UninstallOutputContent struct {
	RecodeAlreadyUninstalled  bool
	SuccessMessage            string
	AlreadyUninstalledMessage string
}

type UninstallOutputHandler interface {
	HandleOutput(UninstallOutput) error
}

type UninstallFeature struct {
	stepper             stepper.Stepper
	outputHandler       UninstallOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewUninstallFeature(
	stepper stepper.Stepper,
	outputHandler UninstallOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) UninstallFeature {

	return UninstallFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (u UninstallFeature) Execute(input UninstallInput) error {
	handleError := func(err error) error {
		u.outputHandler.HandleOutput(UninstallOutput{
			Stepper: u.stepper,
			Error:   err,
		})

		return err
	}

	u.stepper.StartTemporaryStep("Uninstalling Recode")

	cloudService, err := u.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	recodeConfig, err := cloudService.LookupRecodeConfig(
		u.stepper,
	)

	if err != nil {
		if errors.Is(err, entities.ErrRecodeNotInstalled) {
			return u.outputHandler.HandleOutput(UninstallOutput{
				Stepper: u.stepper,
				Content: &UninstallOutputContent{
					RecodeAlreadyUninstalled:  true,
					SuccessMessage:            input.SuccessMessage,
					AlreadyUninstalledMessage: input.AlreadyUninstalledMessage,
				},
			})
		}

		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := recodeConfig.GetCluster(clusterName)

	// In case of error the recode config storage
	// could be created but without cluster
	if err != nil && !errors.As(err, &entities.ErrClusterNotExists{}) {
		return handleError(err)
	}

	if cluster != nil {
		nbOfDevEnvsInCluster, err := recodeConfig.CountDevEnvsInCluster(clusterName)

		if err != nil {
			return handleError(err)
		}

		if nbOfDevEnvsInCluster > 0 {
			return handleError(entities.ErrUninstallExistingDevEnvs)
		}

		err = actions.RemoveCluster(
			u.stepper,
			cloudService,
			recodeConfig,
			cluster,
		)

		if err != nil {
			return handleError(err)
		}
	}

	err = cloudService.RemoveRecodeConfigStorage(
		u.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	return u.outputHandler.HandleOutput(UninstallOutput{
		Stepper: u.stepper,
		Content: &UninstallOutputContent{
			RecodeAlreadyUninstalled:  false,
			SuccessMessage:            input.SuccessMessage,
			AlreadyUninstalledMessage: input.AlreadyUninstalledMessage,
		},
	})
}
