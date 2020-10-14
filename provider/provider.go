package provider

import (
	config "github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/entity"
)

// Provider describes the interface that must be used for each
// repository provider (GitHub, GitLab, BitBucket, etc)
type Provider interface {
	// GetRepos retrieves all the repos that are accessible with
	// the given credentials form the provider.
	GetRepos() []*entity.Repository
}

// NewProvider returns appropriate provider for given name
func NewProvider(c config.Config) Provider {
	if c.ProviderName == "github.com" {
		return NewGithubProvider(c)
	}
	panic(c.ProviderName + " not implemented yet")
}
