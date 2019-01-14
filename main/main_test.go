package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "."
)

func TestParse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parse Suite")
}

var _ = Describe("parseClusterName", func() {
	It("parses a cluster name from a cluster ARN", func() {
		const validARN = "arn:aws:ecs:us-west-2:123456789012:cluster/test-cluster"
		name, err := ParseClusterName(validARN)
		Ω(name).Should(Equal("test-cluster"))
		Ω(err).Should(BeNil())
	})

	It("returns an error when parsing fails", func() {
		const invalidARN = "this:arn:isnt:valid/cluster-name"
		_, err := ParseClusterName(invalidARN)
		Ω(err).ShouldNot(BeNil())
	})
})
