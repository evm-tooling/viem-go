// viemgen generates Go bindings from Ethereum contract ABIs.
//
// Usage:
//
//	viemgen --abi ./MyContract.json --pkg mycontract --out ./contracts/mycontract/
//	viemgen --abi ./MyContract.json --pkg mycontract --name MyContract --out ./contracts/mycontract/
//
// Flags:
//
//	--abi    Path to the ABI JSON file (required)
//	--pkg    Go package name for the generated code (required)
//	--name   Contract name (optional, defaults to package name with first letter capitalized)
//	--out    Output directory (optional, defaults to current directory)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ChefBingbong/viem-go/codegen"
)

func main() {
	var (
		abiPath      string
		packageName  string
		contractName string
		outDir       string
		help         bool
	)

	flag.StringVar(&abiPath, "abi", "", "Path to the ABI JSON file")
	flag.StringVar(&packageName, "pkg", "", "Go package name for the generated code")
	flag.StringVar(&contractName, "name", "", "Contract name (optional)")
	flag.StringVar(&outDir, "out", ".", "Output directory")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&help, "help", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "viemgen - Generate Go bindings from Ethereum contract ABIs\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  viemgen --abi <path> --pkg <package> [--name <name>] [--out <dir>]\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  viemgen --abi ./ERC20.json --pkg erc20 --out ./contracts/erc20/\n")
		fmt.Fprintf(os.Stderr, "  viemgen --abi ./MyContract.json --pkg mycontract --name MyContract\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate required flags
	if abiPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --abi flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if packageName == "" {
		fmt.Fprintln(os.Stderr, "Error: --pkg flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Default contract name to capitalized package name
	if contractName == "" {
		contractName = capitalize(packageName)
	}

	// Read ABI file
	abiJSON, err := os.ReadFile(abiPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading ABI file: %v\n", err)
		os.Exit(1)
	}

	// Create generator
	gen, err := codegen.NewGenerator(packageName, contractName, abiJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating generator: %v\n", err)
		os.Exit(1)
	}

	// Generate code
	code, err := gen.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Create output directory if needed
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Write output file
	outFile := filepath.Join(outDir, packageName+".go")
	if err := os.WriteFile(outFile, code, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", outFile)
}

// capitalize returns the string with the first letter capitalized.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	// Handle special cases like erc20 -> ERC20
	upper := strings.ToUpper(s)
	if upper == "ERC20" || upper == "ERC721" || upper == "ERC1155" {
		return upper
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
