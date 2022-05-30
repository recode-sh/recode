package entities

import "errors"

func (c *Config) SetDevEnv(clusterName string, devEnv *DevEnv) error {
	if devEnv == nil {
		return errors.New("passed dev env is nil")
	}

	c.Clusters[clusterName].DevEnvs[devEnv.Name] = devEnv

	return nil
}

func (c *Config) RemoveDevEnv(clusterName, devEnvName string) error {
	if !c.DevEnvExists(clusterName, devEnvName) {
		return ErrDevEnvNotExists{
			ClusterName: clusterName,
			DevEnvName:  devEnvName,
		}
	}

	delete(c.Clusters[clusterName].DevEnvs, devEnvName)

	return nil
}

func (c *Config) DevEnvExists(clusterName, devEnvName string) bool {
	if !c.ClusterExists(clusterName) {
		return false
	}

	_, devEnvExists := c.Clusters[clusterName].DevEnvs[devEnvName]

	return devEnvExists
}

func (c *Config) GetDevEnv(clusterName, devEnvName string) (*DevEnv, error) {
	if !c.DevEnvExists(clusterName, devEnvName) {
		return nil, ErrDevEnvNotExists{
			ClusterName: clusterName,
			DevEnvName:  devEnvName,
		}
	}

	devEnv := c.Clusters[clusterName].DevEnvs[devEnvName]

	return devEnv, nil
}

func (c *Config) CountDevEnvsInCluster(clusterName string) (int, error) {
	if !c.ClusterExists(clusterName) {
		return 0, ErrClusterNotExists{
			ClusterName: clusterName,
		}
	}

	return len(c.Clusters[clusterName].DevEnvs), nil
}
