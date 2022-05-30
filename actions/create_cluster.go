package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func CreateCluser(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	createClusterErr := cloudService.CreateCluster(
		stepper,
		recodeConfig,
		cluster,
	)

	// "createCLusterErr" is not handled first
	// in order to be able to save partial infrastructure
	err := UpdateClusterInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	if createClusterErr != nil {
		return createClusterErr
	}

	cluster.Status = entities.ClusterStatusCreated
	return UpdateClusterInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
	)
}
