package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKzg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kzg Suite")
}
