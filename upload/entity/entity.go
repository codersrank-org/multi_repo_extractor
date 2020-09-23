package entity

// UploadResult is the result of single repo upload
type UploadResult struct {
	Token string `json:"token"`
}

// MultiUpload is the request body
type MultiUpload struct {
	Results []UploadResultWithRepoName `json:"results"`
}

// UploadResultWithRepoName token-reponame pair
type UploadResultWithRepoName struct {
	Token    string `json:"token"`
	Reponame string `json:"reponame"`
}
