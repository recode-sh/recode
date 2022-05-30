package entities

const (
	DevEnvRepositoryConfigDirectory      = ".recode"
	DevEnvRepositoryDockerfileFileName   = "dev_env.Dockerfile"
	DevEnvRepositoryConfigHooksDirectory = "hooks"
	DevEnvRepositoryInitHookFileName     = "init.sh"
)

type DevEnvRepositoryGitURL string

type ResolvedDevEnvRepository struct {
	Name          string                 `json:"name"`
	Owner         string                 `json:"owner"`
	ExplicitOwner bool                   `json:"explicit_owner"`
	GitURL        DevEnvRepositoryGitURL `json:"git_url"`
	GitHTTPURL    DevEnvRepositoryGitURL `json:"git_http_url"`
}
