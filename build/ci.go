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

// doLint runs the lint target from Makefile (which includes formatting).
func doLint() {
	fmt.Println("==> Running lint (includes gofmt and goimports)...")
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

// doBuild builds all packages.
func doBuild() {
	fmt.Println("==> Building...")
	mustRun("go", "build", "./...")
	fmt.Println("==> Build passed!")
}
