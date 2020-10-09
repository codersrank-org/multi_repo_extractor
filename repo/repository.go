package repo

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	config "github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/repo/entity"
	"github.com/go-git/go-git/v5"
)

// RepositoryService handles repository operations like cloning, updating and processing repos
type RepositoryService interface {
	ProcessRepos(repos []*entity.Repository) []*entity.Repository
	GetTotalRepos() int
	GetRemainingRepos() int
	GetCurrentRepo() *entity.Repository
}

type repositoryService struct {
	RepoInfoExtractorPath string
	RepoInfoExtractorURL  string
	ProviderName          string
	RepoVisibility        string
	Token                 string
	Emails                []string
	SaveRepoPath          string
	AppPath               string
	TotalRepos            int
	ProcessedRepos        int
	CurrentRepository     *entity.Repository
}

// NewRepositoryService constructor
func NewRepositoryService(c config.Config) RepositoryService {
	saveRepoPath := getSaveRepoPath(c.AppPath)
	repositoryService := &repositoryService{
		RepoInfoExtractorPath: c.RepoInfoExtractorPath,
		RepoInfoExtractorURL:  "https://github.com/codersrank-org/repo_info_extractor",
		ProviderName:          c.ProviderName,
		RepoVisibility:        c.RepoVisibility,
		Token:                 c.Token,
		Emails:                c.Emails,
		SaveRepoPath:          saveRepoPath,
		AppPath:               c.AppPath,
	}
	repositoryService.initRepoInfoExtractor()
	return repositoryService
}

func (r *repositoryService) GetTotalRepos() int {
	return r.TotalRepos
}

func (r *repositoryService) GetRemainingRepos() int {
	return r.TotalRepos - r.ProcessedRepos
}

func (r *repositoryService) GetCurrentRepo() *entity.Repository {
	return r.CurrentRepository
}

func (r *repositoryService) ProcessRepos(repos []*entity.Repository) []*entity.Repository {
	r.TotalRepos = len(repos)
	processedRepos := make([]*entity.Repository, 0, len(repos))
	for _, repo := range repos {
		r.ProcessedRepos++
		r.CurrentRepository = repo

		err := r.clone(repo)
		if err != nil {
			continue
		}
		err = r.process(repo)
		if err != nil {
			log.Printf("Couldn't process repo, skipping: %s, error: %s", repo.FullName, err.Error())
			continue
		}
		processedRepos = append(processedRepos, repo)
	}
	return processedRepos
}

func (r *repositoryService) initRepoInfoExtractor() {
	err := cloneRepository(r.RepoInfoExtractorURL, r.RepoInfoExtractorPath, "Repo Info Extractor")
	if err != nil {
		log.Fatalf("Couldn't clone repo_info_extractor: %s", err.Error())
	}
}

func (r *repositoryService) clone(repo *entity.Repository) error {
	// Username is not important, we can use anything as long as it's not an empty string
	repoURL := fmt.Sprintf("https://%s:%s@%s/%s", "username", r.Token, r.ProviderName, repo.FullName)
	repoPath := r.SaveRepoPath + "/" + repo.FullName
	err := cloneRepository(repoURL, repoPath, repo.FullName)
	if err != nil {
		log.Printf("Couldn't clone/update repo, skipping: %s, error: %s", repo.FullName, err.Error())
	}
	return nil
}

func (r *repositoryService) process(repo *entity.Repository) error {
	log.Printf("Processing %s", repo.FullName)

	scriptPath := r.getScriptPath()
	repoPath := r.SaveRepoPath + "/" + repo.FullName

	// Need to chdir to execute scripts because of docker
	os.Chdir(r.RepoInfoExtractorPath)
	cmd := exec.Command(scriptPath, repoPath, "--email="+strings.Join(r.Emails, ","), "--skip_upload", "--headless")

	// We can use these to print repo_info_extractor output to the screen.
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	// Move result to results folder
	sourceLocation := r.RepoInfoExtractorPath + "/repo_data.json.zip"
	targetLocation := getSaveResultPath(r.AppPath) + "/" + strconv.Itoa(repo.ID) + ".zip"

	err = os.Rename(sourceLocation, targetLocation)
	if err != nil {
		return err
	}

	return nil
}

// TODO handle windows (.bat files)
func (r *repositoryService) getScriptPath() string {
	return r.RepoInfoExtractorPath + "/run-docker-headless.sh"
}

// Clone repository from given url to given path
func cloneRepository(url, path, name string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Cloning %s", name)
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout, // TODO add verbose flag to show/hide these.
		})
		if err != nil {
			return err
		}
	} else {
		// If exists, pull latest changes
		log.Printf("Pulling latest changes for %s", name)
		repo, err := git.PlainOpen(path)
		if err != nil {
			return err
		}
		workTree, err := repo.Worktree()
		if err != nil {
			return err
		}
		err = workTree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && !strings.Contains(err.Error(), "already up-to-date") && !strings.Contains(err.Error(), "worktree contains unstaged changes") {
			return err
		}
	}
	return nil
}

func getAppPath() string {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return appPath
}

func getSaveRepoPath(appPath string) string {
	tmpPath := appPath + "/tmp"
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.Mkdir(tmpPath, 0700)
	}
	return tmpPath
}

func getSaveResultPath(appPath string) string {
	resultPath := appPath + "/results"
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		os.Mkdir(resultPath, 0700)
	}
	return resultPath
}
