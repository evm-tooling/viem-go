package abi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestABI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ABI Suite")
}
