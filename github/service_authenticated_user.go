package github

import (
	"context"
	"errors"
)

type AuthenticatedUser struct {
	PrimaryEmail string
	Username     string
	FullName     string
}

func (s Service) GetAuthenticatedUser(
	accessToken string,
) (*AuthenticatedUser, error) {

	client := s.buildClient(accessToken)

	var primaryEmail string
	var getPrimaryEmailErr error
	var getPrimaryEmailChan = make(chan struct{})

	go func() {
		primaryEmail, getPrimaryEmailErr = s.getAuthenticatedUserPrimaryEmail(accessToken)

		close(getPrimaryEmailChan)
	}()

	user, _, getUserErr := client.Users.Get(context.TODO(), "")

	<-getPrimaryEmailChan

	if getUserErr != nil {
		return nil, getUserErr
	}

	if getPrimaryEmailErr != nil {
		return nil, getPrimaryEmailErr
	}

	return &AuthenticatedUser{
		Username:     *user.Login,
		FullName:     *user.Name,
		PrimaryEmail: primaryEmail,
	}, nil
}

func (s Service) getAuthenticatedUserPrimaryEmail(
	accessToken string,
) (string, error) {

	client := s.buildClient(accessToken)

	emails, _, err := client.Users.ListEmails(context.TODO(), nil)

	if err != nil {
		return "", err
	}

	for _, email := range emails {
		if *email.Primary {
			return *email.Email, nil
		}
	}

	return "", errors.New("no primary email found")
}
