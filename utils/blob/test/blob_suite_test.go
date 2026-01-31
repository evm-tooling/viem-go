package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBlob(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Blob Suite")
}
