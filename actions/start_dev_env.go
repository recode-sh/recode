package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func StartDevEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	devEnv.Status = entities.DevEnvStatusStarting
	err := UpdateDevEnvInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
		devEnv,
	)

	if err != nil {
		return err
	}

	startDevEnvErr := cloudService.StartDevEnv(
		stepper,
		recodeConfig,
		cluster,
		devEnv,
	)

	// "startDevEnvErr" is not handled first
	// in order to be able to save partial infrastructure
	err = UpdateDevEnvInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
		devEnv,
	)

	if err != nil {
		return err
	}

	return startDevEnvErr
}
