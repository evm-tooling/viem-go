.PHONY: fmt lint test test-cover verify check release bench

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
