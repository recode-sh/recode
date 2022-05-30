package entities

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

const (
	DevEnvDockerfilesVSCodeExtLabelKey = "sh.recode.vscode.extensions"
	DevEnvDockerfilesReposLabelKey     = "sh.recode.repositories"
	DevEnvRootUser                     = "recode"
)

var (
	DevEnvDockerfilesVSCodeExtLabelSepRegExp = regexp.MustCompile(`\s*,\s*`)
	DevEnvDockerfilesReposLabelSepRegExp     = regexp.MustCompile(`\s*,\s*`)
)

type DevEnvSSHHostKey struct {
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
}

type DevEnvStatus string

const (
	DevEnvStatusCreating DevEnvStatus = "creating"
	DevEnvStatusStarting DevEnvStatus = "starting"
	DevEnvStatusStarted  DevEnvStatus = "started"
	DevEnvStatusStopping DevEnvStatus = "stopping"
	DevEnvStatusStopped  DevEnvStatus = "stopped"
	DevEnvStatusRemoving DevEnvStatus = "removing"
)

type DevEnv struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	InfrastructureJSON       string                   `json:"infrastructure_json"`
	InstanceType             string                   `json:"instance_type"`
	InstancePublicIPAddress  string                   `json:"instance_public_ip_address"`
	InstancePublicHostname   string                   `json:"instance_public_hostname"`
	SSHHostKeys              []DevEnvSSHHostKey       `json:"ssh_host_keys"`
	SSHKeyPairPEMContent     string                   `json:"ssh_key_pair_pem_content"`
	ResolvedUserConfig       ResolvedDevEnvUserConfig `json:"resolved_user_config"`
	ResolvedRepository       ResolvedDevEnvRepository `json:"resolved_repository"`
	Status                   DevEnvStatus             `json:"status"`
	AdditionalPropertiesJSON string                   `json:"additional_properties_json"`
	CreatedAtTimestamp       int64                    `json:"created_at_timestamp"`
}

func NewDevEnv(
	devEnvName string,
	instanceType string,
	resolvedUserConfig ResolvedDevEnvUserConfig,
	resolvedRepository ResolvedDevEnvRepository,
) *DevEnv {

	return &DevEnv{
		ID:                 uuid.NewString(),
		Name:               devEnvName,
		InstanceType:       instanceType,
		SSHHostKeys:        []DevEnvSSHHostKey{},
		ResolvedUserConfig: resolvedUserConfig,
		ResolvedRepository: resolvedRepository,
		Status:             DevEnvStatusCreating,
		CreatedAtTimestamp: time.Now().Unix(),
	}
}

func (d *DevEnv) GetNameSlug() string {
	return BuildDevEnvNameSlug(d.Name)
}

func (d *DevEnv) GetSSHKeyPairName() string {
	return "recode-" + d.GetNameSlug() + "-key-pair"
}

func (d *DevEnv) SetInfrastructureJSON(infrastructure interface{}) error {
	infrastructureJSON, err := json.Marshal(infrastructure)

	if err != nil {
		return err
	}

	d.InfrastructureJSON = string(infrastructureJSON)

	return nil
}

func (d *DevEnv) SetAdditionalPropertiesJSON(additionalProperties interface{}) error {
	additionalPropsJSON, err := json.Marshal(additionalProperties)

	if err != nil {
		return err
	}

	d.AdditionalPropertiesJSON = string(additionalPropsJSON)

	return nil
}

func BuildDevEnvNameSlug(name string) string {
	return slug.Make(name)
}

func BuildDevEnvNameFromResolvedRepo(
	resolvedRepo ResolvedDevEnvRepository,
) string {

	return resolvedRepo.Owner + "/" + resolvedRepo.Name
}

func ParseSSHHostKeys(hostKeysContent string) ([]DevEnvSSHHostKey, error) {
	sanitizedHostKeysContent := strings.TrimSpace(hostKeysContent)
	hostKeys := strings.Split(sanitizedHostKeysContent, "\n")

	parsedHostKeys := []DevEnvSSHHostKey{}

	for _, hostKey := range hostKeys {
		sanitizedHostKey := strings.TrimSpace(hostKey)
		hostKeyComponents := strings.Split(sanitizedHostKey, " ")

		// eg: (ssh-rsa) (AAAAB3NzaC1yc===) (root@ip-10-0-0-200)
		if len(hostKeyComponents) != 3 {
			return nil, fmt.Errorf("invalid host key (\"%s\")", hostKey)
		}

		parsedHostKeys = append(parsedHostKeys, DevEnvSSHHostKey{
			Algorithm:   hostKeyComponents[0],
			Fingerprint: hostKeyComponents[1],
		})
	}

	return parsedHostKeys, nil
}
