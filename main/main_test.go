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

var _ = Describe("ParseClusterName", func() {
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

var _ = Describe("DetailsToMap", func() {
	details := []map[string]string{
		map[string]string{
			"name":  "subnetId",
			"value": "subnet-3afa42b0",
		},
		map[string]string{
			"name":  "networkInterfaceId",
			"value": "eni-d681f702",
		},
		map[string]string{
			"name":  "macAddress",
			"value": "d2:6b:3d:96:1a:68",
		},
		map[string]string{
			"name":  "privateIPv4Address",
			"value": "10.0.0.1",
		},
	}

	It("turns a details array into a map", func() {
		Ω(DetailsToMap(details)).Should(Equal(map[string]string{
			"subnetId":           "subnet-3afa42b0",
			"networkInterfaceId": "eni-d681f702",
			"macAddress":         "d2:6b:3d:96:1a:68",
			"privateIPv4Address": "10.0.0.1",
		}))
	})
})

var _ = Describe("PluckNetworkInterfaceID", func() {
	It("retrieves the network interface id from attachments", func() {
		attachments := []map[string]interface{}{
			map[string]interface{}{
				"id":   "446feb45-f2ca-4481-8e40-9aef61a8f3e1",
				"type": "something-else",
			},
			map[string]interface{}{
				"id":     "13af2a8c-1aff-4a97-b15f-d0cbd78c667d",
				"type":   "eni",
				"status": "ATTACHED",
				"details": []map[string]string{
					map[string]string{
						"name":  "subnetId",
						"value": "subnet-3afa42b0",
					},
					map[string]string{
						"name":  "networkInterfaceId",
						"value": "eni-d681f702",
					},
					map[string]string{
						"name":  "macAddress",
						"value": "d2:6b:3d:96:1a:68",
					},
					map[string]string{
						"name":  "privateIPv4Address",
						"value": "10.0.0.1",
					},
				},
			},
		}

		name, err := PluckNetworkInterfaceID(attachments)
		Ω(name).Should(Equal("eni-d681f702"))
		Ω(err).Should(BeNil())
	})

	It("returns an error when network interface id is missing", func() {
		attachments := []map[string]interface{}{
			map[string]interface{}{
				"id":      "13af2a8c-1aff-4a97-b15f-d0cbd78c667d",
				"type":    "eni",
				"status":  "ATTACHED",
				"details": []map[string]string{},
			},
		}

		_, err := PluckNetworkInterfaceID(attachments)
		Ω(err).ShouldNot(BeNil())
	})

	It("returns an error when eni attachment is missing", func() {
		attachments := []map[string]interface{}{
			map[string]interface{}{
				"id":   "446feb45-f2ca-4481-8e40-9aef61a8f3e1",
				"type": "something-else",
			},
		}

		_, err := PluckNetworkInterfaceID(attachments)
		Ω(err).ShouldNot(BeNil())
	})
})
