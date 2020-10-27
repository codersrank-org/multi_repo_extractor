package provider_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/provider"
)

var _ = Describe("Bitbucket", func() {

	p := provider.NewProvider(config.Config{
		ProviderName:   "bitbucket.org",
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
			httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories?q=is_private+%3D+false&role=contributor", httpmock.NewStringResponder(200, string(getResponseFromFile("../test_fixtures/provider/bitbucket_public.json"))))
			repos := p.GetRepos()
			Expect(len(repos)).To(Equal(10))
			Expect(repos[0].FullName).To(Equal("opensymphony/xwork"))
			Expect(repos[0].Name).To(Equal("xwork"))
			Expect(repos[0].ID).To(Equal("{3f630668-75f1-4903-ae5e-8ea37437e09e}"))
			httpmock.DeactivateAndReset()
		})
	})

})
