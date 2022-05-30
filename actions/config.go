package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func UpdateClusterInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	err := recodeConfig.SetCluster(cluster)

	if err != nil {
		return err
	}

	return cloudService.SaveRecodeConfig(
		stepper,
		recodeConfig,
	)
}

func RemoveClusterInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	err := recodeConfig.RemoveCluster(cluster.Name)

	if err != nil {
		return err
	}

	return cloudService.SaveRecodeConfig(
		stepper,
		recodeConfig,
	)
}

func UpdateDevEnvInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	err := recodeConfig.SetDevEnv(cluster.Name, devEnv)

	if err != nil {
		return err
	}

	return cloudService.SaveRecodeConfig(
		stepper,
		recodeConfig,
	)
}

func RemoveDevEnvInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	err := recodeConfig.RemoveDevEnv(cluster.Name, devEnv.Name)

	if err != nil {
		return err
	}

	return cloudService.SaveRecodeConfig(
		stepper,
		recodeConfig,
	)
}
