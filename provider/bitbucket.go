package provider

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	config "github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/entity"
)

// BitbucketProvider Bitbucket provider used for handling Bitbucket API operations
type BitbucketProvider struct {
	Scheme     string
	BaseURL    string
	Path       string
	Username   string
	Token      string
	Visibility string
}

// NewBitbucketProvider constructor
func NewBitbucketProvider(c config.Config) *BitbucketProvider {
	return &BitbucketProvider{
		Scheme:     "https",
		BaseURL:    "api.bitbucket.org",
		Path:       "2.0/repositories",
		Username:   c.Username,
		Token:      c.Token,
		Visibility: c.RepoVisibility,
	}
}

// GetRepos returns list of repositories with given token and visibility from provider
func (p *BitbucketProvider) GetRepos() []*entity.Repository {
	requestURL := url.URL{
		Scheme: p.Scheme,
		Host:   p.BaseURL,
		Path:   p.Path,
	}

	query := requestURL.Query()
	// role is required otherwise we will get all bitbucket repos.
	query.Set("role", "contributor")

	if p.Visibility == "public" {
		// By default Bitbucket API returns all repositories
		query.Set("q", "is_private = false")
	}

	requestURL.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	request.SetBasicAuth(p.Username, p.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var bitbucketRepos *BitbucketRepository
	err = json.Unmarshal([]byte(body), &bitbucketRepos)
	if err != nil {
		log.Fatal(err)
	}

	repos := make([]*entity.Repository, len(bitbucketRepos.Values))
	for index, repo := range bitbucketRepos.Values {
		repos[index] = &entity.Repository{
			ID:       repo.UUID,
			FullName: repo.FullName,
			Name:     repo.Name,
		}
	}

	return repos
}

// BitbucketRepository response from Bitbucket API
type BitbucketRepository struct {
	Values []struct {
		UUID     string `json:"uuid"`
		FullName string `json:"full_name"`
		Name     string `json:"name"`
	} `json:"values"`
}
