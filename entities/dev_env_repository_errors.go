package entities

type ErrDevEnvRepositoryNotFound struct {
	RepoOwner string
	RepoName  string
}

func (ErrDevEnvRepositoryNotFound) Error() string {
	return "ErrDevEnvRepositoryNotFound"
}
