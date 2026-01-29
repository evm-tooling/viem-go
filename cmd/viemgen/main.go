// viemgen generates Go bindings from Ethereum contract ABIs.
//
// Usage:
//
//	viemgen --abi ./MyContract.json --pkg mycontract
//	viemgen --pkg mycontract                           # Uses default ABI path: _contracts_typed/json/mycontract.json
//	viemgen init                                        # Initialize default directory structure
//
// Default Directories:
//
//	_contracts_typed/
//	├── json/                  # Place ABI JSON files here
//	└── contract_templates/    # Generated Go bindings output here
//
// Flags:
//
//	--abi    Path to the ABI JSON file (default: _contracts_typed/json/<pkg>.json)
//	--pkg    Go package name for the generated code (required)
//	--name   Contract name (optional, defaults to package name capitalized)
//	--out    Output directory (default: _contracts_typed/contract_templates/<pkg>/)
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

const (
	// Default directory structure
	defaultBaseDir     = "_contracts_typed"
	defaultJSONDir     = "json"
	defaultTemplateDir = "contract_templates"
)

func main() {
	// Check for init command
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := initDirectories(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing directories: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Initialized directory structure:")
		fmt.Printf("  %s/\n", defaultBaseDir)
		fmt.Printf("  ├── %s/     # Place ABI JSON files here\n", defaultJSONDir)
		fmt.Printf("  └── %s/     # Generated contracts output here\n", defaultTemplateDir)
		os.Exit(0)
	}

	var (
		abiPath      string
		packageName  string
		contractName string
		outDir       string
		help         bool
	)

	flag.StringVar(&abiPath, "abi", "", "Path to the ABI JSON file (default: _contracts_typed/json/<pkg>.json)")
	flag.StringVar(&packageName, "pkg", "", "Go package name for the generated code (required)")
	flag.StringVar(&contractName, "name", "", "Contract name (optional)")
	flag.StringVar(&outDir, "out", "", "Output directory (default: _contracts_typed/contract_templates/<pkg>/)")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&help, "help", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "viemgen - Generate Go bindings from Ethereum contract ABIs\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  viemgen --pkg <package> [--abi <path>] [--name <name>] [--out <dir>]\n")
		fmt.Fprintf(os.Stderr, "  viemgen init                  # Initialize default directory structure\n\n")
		fmt.Fprintf(os.Stderr, "Default Directories:\n")
		fmt.Fprintf(os.Stderr, "  %s/\n", defaultBaseDir)
		fmt.Fprintf(os.Stderr, "  ├── %s/                  # Place ABI JSON files here\n", defaultJSONDir)
		fmt.Fprintf(os.Stderr, "  └── %s/    # Generated contracts output here\n\n", defaultTemplateDir)
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  viemgen init                                    # Setup directories\n")
		fmt.Fprintf(os.Stderr, "  viemgen --pkg erc20                             # Uses default paths\n")
		fmt.Fprintf(os.Stderr, "  viemgen --pkg erc20 --abi ./custom/ERC20.json   # Custom ABI path\n")
		fmt.Fprintf(os.Stderr, "  viemgen --pkg mytoken --out ./contracts/        # Custom output\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate required flags
	if packageName == "" {
		fmt.Fprintln(os.Stderr, "Error: --pkg flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Ensure default directories exist
	if err := ensureDefaultDirectories(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating default directories: %v\n", err)
		os.Exit(1)
	}

	// Apply defaults if not specified
	if abiPath == "" {
		abiPath = filepath.Join(defaultBaseDir, defaultJSONDir, packageName+".json")
		fmt.Printf("Using default ABI path: %s\n", abiPath)
	}

	if outDir == "" {
		outDir = filepath.Join(defaultBaseDir, defaultTemplateDir, packageName)
		fmt.Printf("Using default output directory: %s\n", outDir)
	}

	// Default contract name to capitalized package name
	if contractName == "" {
		contractName = capitalize(packageName)
	}

	// Check if ABI file exists
	if _, err := os.Stat(abiPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: ABI file not found: %s\n", abiPath)
		fmt.Fprintf(os.Stderr, "\nHint: Place your ABI JSON file at: %s\n", filepath.Join(defaultBaseDir, defaultJSONDir, packageName+".json"))
		fmt.Fprintf(os.Stderr, "  Or specify a custom path with: --abi <path>\n")
		os.Exit(1)
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

// initDirectories creates the default directory structure.
func initDirectories() error {
	dirs := []string{
		filepath.Join(defaultBaseDir, defaultJSONDir),
		filepath.Join(defaultBaseDir, defaultTemplateDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	// Create a README in the json directory
	readmePath := filepath.Join(defaultBaseDir, defaultJSONDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		readme := `# ABI JSON Files

Place your contract ABI JSON files here.

## Naming Convention

Name your files to match the package name you want to use:
- ` + "`erc20.json`" + ` → generates package "erc20"
- ` + "`mytoken.json`" + ` → generates package "mytoken"

## Usage

` + "```bash" + `
# Generate bindings using default paths
viemgen --pkg erc20

# This will:
# - Read from: _contracts_typed/json/erc20.json
# - Output to: _contracts_typed/contract_templates/erc20/erc20.go
` + "```" + `
`
		if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
			return fmt.Errorf("failed to create README: %w", err)
		}
	}

	// Create a .gitkeep in templates directory
	gitkeepPath := filepath.Join(defaultBaseDir, defaultTemplateDir, ".gitkeep")
	if _, err := os.Stat(gitkeepPath); os.IsNotExist(err) {
		if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep: %w", err)
		}
	}

	return nil
}

// ensureDefaultDirectories creates the default directories if they don't exist.
func ensureDefaultDirectories() error {
	dirs := []string{
		filepath.Join(defaultBaseDir, defaultJSONDir),
		filepath.Join(defaultBaseDir, defaultTemplateDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
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
