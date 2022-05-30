package github

import (
	"context"

	"github.com/google/go-github/v43/github"
)

func (s Service) CreateGPGKey(
	accessToken string,
	publicKeyContent string,
) (*github.GPGKey, error) {

	client := s.buildClient(accessToken)

	key, _, err := client.Users.CreateGPGKey(
		context.TODO(),
		publicKeyContent,
	)

	return key, err
}

func (s Service) RemoveGPGKey(
	accessToken string,
	gpgKeyID int64,
) error {

	client := s.buildClient(accessToken)

	_, err := client.Users.DeleteGPGKey(
		context.TODO(),
		gpgKeyID,
	)

	return err
}
