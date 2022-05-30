package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func StopDevEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	devEnv.Status = entities.DevEnvStatusStopping
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

	stopDevEnvErr := cloudService.StopDevEnv(
		stepper,
		recodeConfig,
		cluster,
		devEnv,
	)

	// "stopDevEnvErr" is not handled first
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

	return stopDevEnvErr
}
