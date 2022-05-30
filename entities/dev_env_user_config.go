package entities

const (
	DevEnvUserConfigDefaultRepoOwner    = "recode-sh"
	DevEnvUserConfigRepoName            = ".recode"
	DevEnvUserConfigDockerfileFileName  = "dev_env.Dockerfile"
	DevEnvUserConfigDockerfileRootImage = "recodesh/base-dev-env"
	DevEnvUserConfigDockerfileImageName = "user_dev_env"
)

type ResolvedDevEnvUserConfig struct {
	RepoOwner      string                 `json:"repo_owner"`
	RepoName       string                 `json:"repo_name"`
	RepoGitURL     DevEnvRepositoryGitURL `json:"repo_git_url"`
	RepoGitHTTPURL DevEnvRepositoryGitURL `json:"repo_git_http_url"`
}
