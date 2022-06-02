package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func RemoveDevEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
	preRemoveHook entities.HookRunner,
) error {

	devEnv.Status = entities.DevEnvStatusRemoving
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

	removeDevEnvErr := cloudService.RemoveDevEnv(
		stepper,
		recodeConfig,
		cluster,
		devEnv,
	)

	// "removeDevEnvErr" is not handled first
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

	if removeDevEnvErr != nil {
		return removeDevEnvErr
	}

	if preRemoveHook != nil {
		err = preRemoveHook.Run(
			cloudService,
			recodeConfig,
			cluster,
			devEnv,
		)

		if err != nil {
			return err
		}
	}

	return RemoveDevEnvInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
		devEnv,
	)
}
