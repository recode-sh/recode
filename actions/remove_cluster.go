package actions

import (
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/stepper"
)

func RemoveCluster(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	cluster.Status = entities.ClusterStatusRemoving
	err := UpdateClusterInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	removeClusterErr := cloudService.RemoveCluster(
		stepper,
		recodeConfig,
		cluster,
	)

	// "removeClusterErr" is not handled first
	// in order to be able to save partial infrastructure
	err = UpdateClusterInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	if removeClusterErr != nil {
		return removeClusterErr
	}

	return RemoveClusterInConfig(
		stepper,
		cloudService,
		recodeConfig,
		cluster,
	)
}
