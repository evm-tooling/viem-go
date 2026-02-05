.PHONY: fmt lint test test-cover verify check release bench fmt-ts lint-ts check-ts oncommit

# ============================================================================
# Go targets
# ============================================================================

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './benchmarks/*')
	goimports -local github.com/ChefBingbong/viem-go -w $$(find . -name '*.go' -not -path './benchmarks/*')

lint:
	golangci-lint run

verify:
	golangci-lint config verify

test:
	go test -v -race $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... | grep -v '/benchmarks/')

test-cover:
	go test -v -race -coverprofile=coverage.out $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... | grep -v '/benchmarks/')

check: fmt lint test

# ============================================================================
# TypeScript targets (examples/viem workspace)
# ============================================================================

# Format TypeScript files using Biome
fmt-ts:
	@echo "==> Formatting TypeScript files..."
	@cd examples/viem && bun run format

# Lint TypeScript files using Biome
lint-ts:
	@echo "==> Linting TypeScript files..."
	@cd examples/viem && bun run lint

# Fix TypeScript lint issues using Biome
lint-ts-fix:
	@echo "==> Fixing TypeScript lint issues..."
	@cd examples/viem && bun run lint:fix

# Check TypeScript types
check-ts:
	@echo "==> Checking TypeScript types..."
	@cd examples/viem && bun run check

# ============================================================================
# Combined targets
# ============================================================================

# Pre-commit check: verifies formatting and linting for both Go and TypeScript
oncommit: fmt lint lint-ts
	@echo "==> All pre-commit checks passed!"

# Full check for both Go and TypeScript
check-all: check check-ts lint-ts
	@echo "==> All checks passed!"

# Release workflow: creates tag, drafts prerelease, and opens PR to production
# Usage: make release VERSION=v0.0.4
# Requires: gh CLI authenticated with repo write access
release:
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=v0.0.4)
endif
	@echo "==> Checking gh CLI authentication..."
	@gh auth status || (echo "Error: Please run 'gh auth login' first" && exit 1)
	@echo "==> Ensuring we're on master branch..."
	@git checkout master
	@git pull origin master
	@echo "==> Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "==> Creating prerelease on GitHub..."
	@gh release create $(VERSION) \
		--title "$(VERSION)" \
		--generate-notes \
		--prerelease
	@echo "==> Creating PR from master to production..."
	@gh pr create \
		--base production \
		--head master \
		--title "Release $(VERSION)" \
		--body "## Release $(VERSION)$$( echo )$$( echo )This PR releases $(VERSION) to production.$$( echo )$$( echo )When merged, the prerelease will be automatically published as the latest release."
	@echo "==> Done! Review and merge the PR to publish the release."

# Run cross-language benchmarks (viem-go vs viem TypeScript)
# Requires: Foundry (anvil), Node.js, bun (optional for compare)
bench:
	$(MAKE) -C benchmarks bench
