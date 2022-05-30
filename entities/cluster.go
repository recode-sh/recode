package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

const (
	DefaultClusterName = "default"
)

type ClusterStatus string

const (
	ClusterStatusCreating ClusterStatus = "creating"
	ClusterStatusCreated  ClusterStatus = "created"
	ClusterStatusRemoving ClusterStatus = "removing"
)

type Cluster struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	DefaultInstanceType string             `json:"default_instance_type"`
	InfrastructureJSON  string             `json:"infrastructure_json"`
	DevEnvs             map[string]*DevEnv `json:"dev_envs"`
	IsDefault           bool               `json:"is_default"`
	Status              ClusterStatus      `json:"status"`
	CreatedAtTimestamp  int64              `json:"created_at_timestamp"`
}

func NewCluster(
	clusterName string,
	defaultInstanceType string,
	isDefaultCluster bool,
) *Cluster {

	return &Cluster{
		ID:                  uuid.NewString(),
		Name:                clusterName,
		DefaultInstanceType: defaultInstanceType,
		DevEnvs:             map[string]*DevEnv{},
		IsDefault:           isDefaultCluster,
		Status:              ClusterStatusCreating,
		CreatedAtTimestamp:  time.Now().Unix(),
	}
}

func (c *Cluster) GetNameSlug() string {
	return slug.Make(c.Name)
}

func (c *Cluster) SetInfrastructureJSON(infrastructure interface{}) error {
	infrastructureJSON, err := json.Marshal(infrastructure)

	if err != nil {
		return err
	}

	c.InfrastructureJSON = string(infrastructureJSON)

	return nil
}
