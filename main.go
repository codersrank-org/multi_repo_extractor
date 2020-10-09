package main

import (
	"github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/provider"
	repo "github.com/codersrank-org/multi_repo_repo_extractor/repo"
	upload "github.com/codersrank-org/multi_repo_repo_extractor/upload"
)

func main() {

	// TODO implement auto-update (versioning etc.)

	config := config.ParseFlags()

	provider := provider.NewProvider(config)
	repositoryService := repo.NewRepositoryService(config)
	codersrankService := upload.NewCodersrankService(config)

	repos := provider.GetRepos()
	processedRepos := repositoryService.ProcessRepos(repos)
	codersrankService.UploadRepos(processedRepos)
}
