package github

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/recode-sh/recode/entities"
	giturls "github.com/whilp/git-urls"
)

func BuildGitHTTPURL(repoOwner, repoName string) entities.DevEnvRepositoryGitURL {
	return entities.DevEnvRepositoryGitURL(fmt.Sprintf(
		"https://github.com/%s/%s.git",
		url.PathEscape(repoOwner),
		url.PathEscape(repoName),
	))
}

func BuildGitURL(repoOwner, repoName string) entities.DevEnvRepositoryGitURL {
	return entities.DevEnvRepositoryGitURL(fmt.Sprintf(
		"git@github.com:%s/%s.git",
		url.PathEscape(repoOwner),
		url.PathEscape(repoName),
	))
}

type ParsedGitHubRepositoryName struct {
	Owner         string
	ExplicitOwner bool
	Name          string
}

func ParseRepositoryName(
	repositoryName string,
	defaultRepositoryOwner string,
) (*ParsedGitHubRepositoryName, error) {

	errInvalidGitHubURL := errors.New("ErrInvalidGitHubURL")

	// Handle git@github.com:recode-sh/recode.git
	repositoryNameAsURL, err := giturls.Parse(repositoryName)

	if err != nil {
		// Handle https://github.com/recode-sh/recode.git
		repositoryNameAsURL, err = url.Parse(repositoryName)
	}

	// Not an URL (eg: recode) or only path (eg: recode-sh/recode)
	if err != nil || len(repositoryNameAsURL.Hostname()) == 0 {
		repositoryNameParts := strings.Split(repositoryName, "/")

		if len(repositoryNameParts) > 2 {
			return nil, errInvalidGitHubURL
		}

		if len(repositoryNameParts) == 1 { // recode
			return &ParsedGitHubRepositoryName{
				ExplicitOwner: false,
				Owner:         defaultRepositoryOwner,
				Name:          repositoryNameParts[0],
			}, nil
		}

		return &ParsedGitHubRepositoryName{ // recode-sh/recode
			ExplicitOwner: true,
			Owner:         repositoryNameParts[0],
			Name:          repositoryNameParts[1],
		}, nil
	}

	host := repositoryNameAsURL.Hostname()

	if host != "github.com" {
		return nil, errInvalidGitHubURL
	}

	path := strings.TrimPrefix(repositoryNameAsURL.Path, "/")
	pathComponents := strings.Split(path, "/")

	if len(pathComponents) < 2 {
		return nil, errInvalidGitHubURL
	}

	githubRepositoryOwner := pathComponents[0]
	githubRepositoryName := strings.TrimSuffix(pathComponents[1], ".git")

	return &ParsedGitHubRepositoryName{
		ExplicitOwner: true,
		Owner:         githubRepositoryOwner,
		Name:          githubRepositoryName,
	}, nil
}
