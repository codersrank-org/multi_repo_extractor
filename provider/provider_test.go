package provider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/codersrank-org/multi_repo_repo_extractor/config"
	"github.com/codersrank-org/multi_repo_repo_extractor/provider"
)

var _ = Describe("Providers", func() {

	Describe("Creating provider", func() {
		It("should return with correct provider", func() {
			provider := provider.NewProvider(config.Config{
				ProviderName: "github.com",
			})
			Expect(provider).NotTo(BeNil())
		})
	})

})
