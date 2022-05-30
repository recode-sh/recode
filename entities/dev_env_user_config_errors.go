package entities

type ErrInvalidDevEnvUserConfig struct {
	RepoOwner string
	Reason    string
}

func (ErrInvalidDevEnvUserConfig) Error() string {
	return "ErrInvalidDevEnvUserConfig"
}
