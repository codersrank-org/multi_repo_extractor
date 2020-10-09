package provider

import (
	config "github.com/codersrank-org/multi_repo_repo_extractor/config"
	providers "github.com/codersrank-org/multi_repo_repo_extractor/provider/providers"
	"github.com/codersrank-org/multi_repo_repo_extractor/repo/entity"
)

// Provider is where repositories hosted (e.g. github, gitlab etc.)
type Provider interface {
	GetRepos() []*entity.Repository
}

// NewProvider returns appropriate provider for given name
func NewProvider(c config.Config) Provider {
	if c.ProviderName == "github.com" {
		return providers.NewGithubProvider(c)
	}
	panic(c.ProviderName + " not implemented yet")
}
