package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

// ParseFlags parses flags and environment variables
func ParseFlags() Config {

	var provider, emailString, token, repoVisibility string
	flag.StringVar(&provider, "provider", "github.com", "Provider for repos. Only github.com is supported now.")
	flag.StringVar(&token, "token", "", "Token for accessing repositories. You can also set this with TOKEN environment variable.")
	flag.StringVar(&emailString, "emails", "", "Your emails which are used when making the commits. Provide a comma separated list for multiple emails (e.g. \"one@mail.com,two@email.com\")")
	flag.StringVar(&repoVisibility, "repo_visibility", "private", "Which repos do you want to get processed? Options: all, public and private.")

	flag.Parse()

	// After getting flags, check environment variables
	// If there is an env_var and related variable hasn't specified as a flag, we will use it
	// Which means flags override env_vars
	if os.Getenv("TOKEN") != "" && token == "" {
		token = os.Getenv("TOKEN")
	}

	appPath := getAppPath()
	repoInfoExtractorPath := getDefaultRepoInfoExtractorPath(appPath)
	// Set this if you don't want to download repo_info_extractor and want to use your local version
	if os.Getenv("REPO_EXTRACTOR") != "" {
		repoInfoExtractorPath = os.Getenv("REPO_EXTRACTOR")
	}

	if provider != "github.com" {
		log.Fatal("Only supported provider is github.com.")
	}

	// We can add more checks here
	if token == "" {
		log.Fatal("You need to provide a valid token.")
	}

	var emails []string
	if emailString == "" {
		log.Fatal("You need to provide at least one email.")
	} else {
		emails := strings.Split(emailString, ",")
		for i := range emails {
			emails[i] = strings.TrimSpace(emails[i])
		}
	}

	if repoVisibility != "all" && repoVisibility != "public" && repoVisibility != "private" {
		log.Fatal("Valid values for repo_visibility are: all, public and private.")
	}

	return Config{
		ProviderName:          provider,
		Token:                 token,
		Emails:                emails,
		RepoVisibility:        repoVisibility,
		AppPath:               appPath,
		RepoInfoExtractorPath: repoInfoExtractorPath,
	}
}

// Default path is relative to the current directory
func getDefaultRepoInfoExtractorPath(appPath string) string {
	return appPath + "/repo_info_extractor"
}

func getAppPath() string {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return appPath
}

// Config flags and paths
type Config struct {
	ProviderName          string
	Token                 string
	Emails                []string
	RepoVisibility        string
	AppPath               string
	RepoInfoExtractorPath string
}
