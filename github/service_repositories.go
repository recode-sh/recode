package github

import (
	"context"

	"github.com/google/go-github/v43/github"
)

func (s Service) CreateRepository(
	accessToken string,
	organization string,
	properties *github.Repository,
) (*github.Repository, error) {

	client := s.buildClient(accessToken)

	repository, _, err := client.Repositories.Create(
		context.TODO(),
		organization,
		properties,
	)

	return repository, err
}

func (s Service) DoesRepositoryExist(
	accessToken string,
	repositoryOwner string,
	repositoryName string,
) (bool, error) {

	client := s.buildClient(accessToken)

	repository, _, err := client.Repositories.Get(
		context.TODO(),
		repositoryOwner,
		repositoryName,
	)

	if s.IsNotFoundError(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return repository != nil, nil
}

func (s Service) GetFileContentFromRepository(
	accessToken string,
	repositoryOwner string,
	repositoryName string,
	filePath string,
) (string, error) {

	client := s.buildClient(accessToken)

	fileContent, _, _, err := client.Repositories.GetContents(
		context.TODO(),
		repositoryOwner,
		repositoryName,
		filePath,
		nil,
	)

	if err != nil {
		return "", err
	}

	return fileContent.GetContent()
}
