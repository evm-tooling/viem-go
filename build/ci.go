//go:build ignore

/*
The ci command is called from Continuous Integration scripts.

Usage: go run build/ci.go <command>

Available commands are:

	lint    -- runs formatting checks and linters via Makefile
	test    -- runs the test suite
	build   -- builds all packages
*/
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if _, err := os.Stat(filepath.Join("build", "ci.go")); os.IsNotExist(err) {
		log.Fatal("this script must be run from the root of the repository")
	}
	if len(os.Args) < 2 {
		log.Fatal("need subcommand as first argument")
	}

	switch os.Args[1] {
	case "lint":
		doLint()
	case "test":
		doTest()
	case "test-cover":
		doTestCover()
	case "build":
		doBuild()
	default:
		log.Fatal("unknown command ", os.Args[1])
	}
}

// runCommand executes a command and streams output to stdout/stderr.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// mustRun executes a command and exits on failure.
func mustRun(name string, args ...string) {
	if err := runCommand(name, args...); err != nil {
		log.Fatalf("command failed: %s %v: %v", name, args, err)
	}
}

// installTools installs required Go tools if not present.
func installTools() {
	tools := []struct {
		name string
		pkg  string
	}{
		{"goimports", "golang.org/x/tools/cmd/goimports@latest"},
		{"golangci-lint", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0"},
	}

	for _, tool := range tools {
		if _, err := exec.LookPath(tool.name); err != nil {
			fmt.Printf("==> Installing %s...\n", tool.name)
			mustRun("go", "install", tool.pkg)
		}
	}
}

// doLint runs the lint target from Makefile (which includes formatting).
func doLint() {
	fmt.Println("==> Running lint (includes gofmt and goimports)...")
	installTools()
	mustRun("make", "lint")
	fmt.Println("==> Lint passed!")
}

// doTest runs the test suite.
func doTest() {
	fmt.Println("==> Running tests...")
	mustRun("make", "test")
	fmt.Println("==> Tests passed!")
}

// doTestCover runs tests with coverage.
func doTestCover() {
	fmt.Println("==> Running tests with coverage...")
	mustRun("make", "test-cover")
	fmt.Println("==> Tests passed!")
}

// doBuild builds all packages (excluding benchmarks and test-only directories).
func doBuild() {
	fmt.Println("==> Building...")
	// Exclude benchmarks and test directories from build
	// Test directories only contain *_test.go files which can't be built standalone
	cmd := exec.Command("sh", "-c", "go build $(go list ./... | grep -v '/benchmarks/' | grep -v '/test$')")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("build failed: %v", err)
	}
	fmt.Println("==> Build passed!")
}
