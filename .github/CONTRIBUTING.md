# Contributing

Thanks for your interest in contributing to viem-go! Please take a moment to review this document **before submitting a pull request.**

If you want to contribute but aren't sure where to start, you can create a [new discussion](https://github.com/ChefBingbong/viem-go/discussions).

## Rules

1. Significant changes to the API or implementation must be reviewed before a Pull Request is created. Create a [Feature Request](https://github.com/ChefBingbong/viem-go/discussions/new?category=ideas) first to discuss any API changes or new ideas.
2. Contributors must be humans, not bots.
3. First time contributions should not contain only spelling or grammatical fixes.

## Basic Guide

This guide is intended to help you get started with contributing. By following these steps, you will understand the development process and workflow.

1. [Cloning the repository](#cloning-the-repository)
2. [Installing Go](#installing-go)
3. [Installing dependencies](#installing-dependencies)
4. [Running the test suite](#running-the-test-suite)
5. [Code formatting and linting](#code-formatting-and-linting)
6. [Project structure](#project-structure)
7. [Submitting a pull request](#submitting-a-pull-request)
8. [Versioning](#versioning)

---

### Cloning the repository

To start contributing to the project, clone it to your local machine using git:

```bash
git clone https://github.com/ChefBingbong/viem-go.git
```

Or the [GitHub CLI](https://cli.github.com):

```bash
gh repo clone ChefBingbong/viem-go
```

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Installing Go

viem-go requires **Go 1.24 or higher**. You can check your Go version with:

```bash
go version
```

If you need to install or update Go, download it from the [official website](https://go.dev/dl/) or use a version manager like [gvm](https://github.com/moovweb/gvm) or [asdf](https://github.com/asdf-vm/asdf).

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Installing dependencies

Once in the project's root directory, run the following command to download dependencies:

```bash
go mod download
```

You'll also need to install the development tools for linting and formatting:

```bash
# Install goimports (for import organization)
go install golang.org/x/tools/cmd/goimports@latest

# Install golangci-lint (for linting)
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
```

> **Note:** The CI script (`go run build/ci.go lint`) will automatically install these tools if they're not present.

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Running the test suite

Run the test suite with:

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run all checks (format, lint, test)
make check
```

Or using the CI script:

```bash
go run build/ci.go test
go run build/ci.go test-cover
```

When adding new features or fixing bugs, it's important to add test cases to cover the new/updated behavior. Tests are located in `test/` subdirectories within each package (e.g., `abi/test/`, `client/test/`).

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Code formatting and linting

viem-go uses `gofmt`, `goimports`, and `golangci-lint` for code quality.

```bash
# Format code
make fmt

# Run linter
make lint

# Verify linter config
make verify

# Run all checks (format + lint + test)
make check
```

The project uses a pre-push hook that runs lint and tests before pushing. To set up the hook:

```bash
cp .github/hooks/pre-push .git/hooks/pre-push
chmod +x .git/hooks/pre-push
```

#### Linter Configuration

The linter configuration is in `.golangci.yml`. The following linters are enabled:

- `errcheck` - Check for unchecked errors
- `govet` - Reports suspicious constructs
- `ineffassign` - Detects ineffectual assignments
- `staticcheck` - Static analysis
- `unused` - Finds unused code
- `misspell` - Finds misspellings
- `unconvert` - Finds unnecessary conversions
- `bodyclose` - Checks HTTP response body is closed

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Project structure

```
viem-go/
├── abi/            # ABI encoding/decoding utilities
├── accounts/       # Account management and signing
├── chain/          # Chain definitions and utilities
├── client/         # Client implementations (public, wallet)
│   ├── decorators/ # Client decorators
│   └── transport/  # HTTP, WebSocket, fallback transports
├── codegen/        # Code generation utilities
├── contract/       # Contract interaction utilities
├── contracts/      # Pre-built contract bindings (ERC20, ERC721, ERC1155)
├── crypto/         # Cryptographic utilities (hashing, signatures)
├── types/          # Core type definitions
├── utils/          # Utility functions (unit conversion, hex, etc.)
├── build/          # Build and CI scripts
├── cmd/            # CLI tools (viemgen)
├── examples/       # Example code
└── site/           # Documentation site
```

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Submitting a pull request

When you're ready to submit a pull request, follow these guidelines:

1. **Create a feature branch** from `master`:
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes** and ensure all tests pass:
   ```bash
   make check
   ```

3. **Commit your changes** with a clear, descriptive message:
   - Use the [Imperative Mood](https://en.wikipedia.org/wiki/Imperative_mood) (e.g., `Add something`, `Fix something`)
   - Keep the subject line under 72 characters
   - Reference related issues (e.g., `Fixes #123`)

4. **Push your branch** and open a pull request against `master`.

When you submit a pull request, GitHub will automatically lint, build, and test your changes. If you see an ❌, inspect the logs through the GitHub UI to find the cause.

#### Pull Request Checklist

Before submitting, ensure:

- [ ] Code follows the project's style guidelines (`make fmt` and `make lint` pass)
- [ ] All tests pass (`make test`)
- [ ] New functionality includes appropriate tests
- [ ] Documentation is updated if needed
- [ ] Commit messages are clear and descriptive

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

### Versioning

viem-go follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner
- **PATCH** version for backwards compatible bug fixes

#### Release Process

1. Changes are developed on feature branches and merged to `master`
2. When ready for release, changes are merged to the `production` branch
3. A tag is created and a prerelease is drafted on GitHub
4. When the PR from `production` to `master` is merged, the prerelease is automatically published as the latest release

> **Note:** Only maintainers can create releases. If you believe your contribution warrants a release, mention it in your PR.

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

## Adding a Chain

If you want to add a new chain definition to `chain/definitions/`, follow these steps:

### Requirements

- **Must have**:
  - A unique Chain ID
  - A human-readable name
  - Native currency information (name, symbol, decimals)
  - At least one public RPC URL
- **Nice to have**:
  - Block explorer URL
  - Multicall3 contract address

### Steps

1. Create a new file in `chain/definitions/` (e.g., `mychain.go`)
2. Define your chain using `DefineChain`:

```go
package definitions

import "github.com/ChefBingbong/viem-go/chain"

var MyChain = chain.DefineChain(chain.ChainConfig{
    ID:   12345,
    Name: "My Chain",
    NativeCurrency: chain.NativeCurrency{
        Name:     "MyCoin",
        Symbol:   "MYC",
        Decimals: 18,
    },
    RpcUrls: chain.RpcUrls{
        Default: chain.RpcUrl{
            HTTP: []string{"https://rpc.mychain.io"},
        },
    },
    BlockExplorers: chain.BlockExplorers{
        Default: chain.BlockExplorer{
            Name: "MyChain Explorer",
            URL:  "https://explorer.mychain.io",
        },
    },
})
```

3. Export it in `chain/definitions/util.go` if needed
4. Add tests to verify the chain configuration

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>

---

## Getting Help

- **Questions**: Open a [Discussion](https://github.com/ChefBingbong/viem-go/discussions/new?category=q-a)
- **Bug Reports**: Open an [Issue](https://github.com/ChefBingbong/viem-go/issues/new?template=bug_report.yml)
- **Feature Requests**: Open a [Discussion](https://github.com/ChefBingbong/viem-go/discussions/new?category=ideas)

---

<div>
  ✅ Now you're ready to contribute to viem-go!
</div>

<div align="right">
  <a href="#basic-guide">&uarr; back to top</a>
</div>
