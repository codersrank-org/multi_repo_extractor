package provider_test

import (
	"io/ioutil"
	"os"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/provider"
)

var _ = Describe("Providers", func() {

	p := provider.NewProvider(config.Config{
		ProviderName:   "github.com",
		Token:          "token",
		RepoVisibility: "public",
	})

	Describe("Creating provider", func() {
		It("should return with correct provider", func() {
			Expect(p).NotTo(BeNil())
		})
	})

	Describe("Getting repositories", func() {
		It("should get repositories of the user", func() {
			httpmock.Activate()
			httpmock.RegisterResponder("GET", "https://api.github.com/user/repos?visibility=public", httpmock.NewStringResponder(200, string(getResponseFromFile("../test_fixtures/provider/github_public.json"))))
			repos := p.GetRepos()
			Expect(len(repos)).To(Equal(20))
			Expect(repos[0].FullName).To(Equal("alimgiray/bdd"))
			Expect(repos[0].Name).To(Equal("bdd"))
			Expect(repos[0].ID).To(Equal("134240628"))
			httpmock.DeactivateAndReset()
		})
	})

})

func getResponseFromFile(filePath string) []byte {
	responseFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer responseFile.Close()

	byteValue, err := ioutil.ReadAll(responseFile)
	if err != nil {
		panic(err)
	}
	return byteValue
}
