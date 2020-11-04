package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var currentMajorVersion = 1
var currentMinorVersion = 0
var currentPatchVersion = 0

// CheckUpdates checks github to see if there is a new version and if there is one, downloads it.
func CheckUpdates() {
	log.Println("Checking for new versions")
	latestVersion, err := getLatestVersion()
	if err != nil {
		log.Printf("Couldn't get latest release from Github, skipping update. Error: %s", err.Error())
		return
	}
	log.Println("latest version", latestVersion)
}

func getLatestVersion() (version, error) {
	url := "https://api.github.com/repos/codersrank-org/multi_repo_extractor/releases/latest"

	v := version{}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return v, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return v, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var r *release
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		return v, err
	}

	// "v" is not part of semantic versioning
	r.Name = strings.TrimLeft(r.Name, "v")

	// Regex for finding Major, Minor and Patch versions
	// Taken from here: https://semver.org/
	regex := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	matches := regex.FindAllString(r.Name, -1)
	if len(matches) == 0 {
		return v, errors.New("Couldn't parse current version")
	}

	matches = strings.Split(matches[0], ".")

	if len(matches) != 3 {
		return v, errors.New("Couldn't parse current version")
	}

	v.Major, _ = strconv.Atoi(matches[0])
	v.Minor, _ = strconv.Atoi(matches[1])
	v.Patch, _ = strconv.Atoi(matches[2])

	return v, err
}

type version struct {
	Major int
	Minor int
	Patch int
}

type release struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string      `json:"url"`
		ID       int         `json:"id"`
		NodeID   string      `json:"node_id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}
