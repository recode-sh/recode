package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func CreateDevEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	createDevEnvErr := cloudService.CreateDevEnv(
		stepper,
		recodeConfig,
		cluster,
		devEnv,
	)

	// "createDevEnvErr" is not handled first
	// in order to be able to save partial infrastructure
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

	return createDevEnvErr
}
