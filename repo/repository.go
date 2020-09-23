package repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/codersrank-org/multi_repo_repo_extractor/repo/entity"
	"github.com/go-git/go-git/v5"
)

// RepositoryService handles repository operations like cloning, updating and processing repos
type RepositoryService interface {
	InitRepoInfoExtractor()
	GetReposFromProvider() []*entity.GithubRepository
	Clone(repo *entity.GithubRepository) error
	Process(repo *entity.GithubRepository) error
}

type repositoryService struct {
	RepoInfoExtractorPath string
	RepoInfoExtractorURL  string
	GithubAPI             string
	Provider              string
	RepoVisibility        string
	Token                 string
	Emails                []string
	SaveRepoPath          string
	AppPath               string
}

// NewRepositoryService constructor
func NewRepositoryService(repoInfoExtractorPath, provider, repoVisibility, token string, emails []string) RepositoryService {
	appPath := getAppPath()
	saveRepoPath := getSaveRepoPath(appPath)
	return &repositoryService{
		RepoInfoExtractorPath: repoInfoExtractorPath,
		RepoInfoExtractorURL:  "https://github.com/codersrank-org/repo_info_extractor",
		GithubAPI:             "https://api.github.com/user/repos",
		Provider:              provider,
		RepoVisibility:        repoVisibility,
		Token:                 token,
		Emails:                emails,
		SaveRepoPath:          saveRepoPath,
		AppPath:               appPath,
	}
}

func (r *repositoryService) InitRepoInfoExtractor() {
	log.Println("Initializing repo_info_extractor..")
	// If not exists, clone
	// If exists just pull latest changes
	if _, err := os.Stat(r.RepoInfoExtractorPath); os.IsNotExist(err) {
		_, err := git.PlainClone(r.RepoInfoExtractorPath, false, &git.CloneOptions{
			URL:      r.RepoInfoExtractorURL,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Fatalf("Couldn't clone repo_info_extractor: %s", err.Error())
		}
	} else {
		repo, err := git.PlainOpen(r.RepoInfoExtractorPath)
		if err != nil {
			log.Fatal(err)
		}
		workTree, err := repo.Worktree()
		if err != nil {
			log.Fatal(err)
		}
		// Having no new commits doesn't block us, so ignore that error
		err = workTree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && !strings.Contains(err.Error(), "already up-to-date") {
			log.Fatal(err)
		}
	}
}

func (r *repositoryService) GetReposFromProvider() []*entity.GithubRepository {
	requestURL := fmt.Sprintf("%s?visibility=%s", r.GithubAPI, r.RepoVisibility)
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.Token))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var repos []*entity.GithubRepository
	err = json.Unmarshal([]byte(body), &repos)
	if err != nil {
		log.Fatal(err)
	}

	return repos
}

func (r *repositoryService) Clone(repo *entity.GithubRepository) error {
	// Username is not important, we can use anything as long as it's not an empty string
	repoURL := fmt.Sprintf("https://%s:%s@%s/%s", "username", r.Token, r.Provider, repo.FullName)
	repoPath := r.SaveRepoPath + "/" + repo.FullName
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Printf("Cloning %s", repo.FullName)
		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout, // TODO add verbose flag to show/hide these.
		})
		if err != nil {
			return err
		}
	} else {
		// If exists, pull latest changes
		log.Printf("Pulling latest changes for %s", repo.FullName)
		repo, err := git.PlainOpen(repoPath)
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

func (r *repositoryService) Process(repo *entity.GithubRepository) error {
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
