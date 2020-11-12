package main

import (
	"github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/entity"
	"github.com/codersrank-org/multi_repo_repo_extractor/provider"
	"github.com/codersrank-org/multi_repo_repo_extractor/repo"
	"github.com/codersrank-org/multi_repo_repo_extractor/upload"
)

func main() {
	config.CheckUpdates()
	config := config.ParseFlags()

	providers := make([]provider.Provider, 1)
	providers[0] = provider.NewProvider(config)
	repositoryService := repo.NewRepositoryService(config)
	codersrankService := upload.NewCodersrankService(config)

	repos := make([]*entity.Repository, 0)
	for _, provider := range providers {
		repos = append(repos, provider.GetRepos()...)
	}
	processedRepos := repositoryService.ProcessRepos(repos)
	codersrankService.UploadRepos(processedRepos)
}
