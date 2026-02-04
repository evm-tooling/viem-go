.PHONY: fmt lint test test-cover verify check release

fmt:
	gofmt -w .
	goimports -local github.com/ChefBingbong/viem-go -w .

lint:
	golangci-lint run

verify:
	golangci-lint config verify

test:
	go test -v -race $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

test-cover:
	go test -v -race -coverprofile=coverage.out $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

check: fmt lint test

# Release workflow: triggers CI to create tag, prerelease, and PR
# Usage: make release VERSION=v0.0.4
# Requires: gh CLI authenticated with repo write access
release:
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=v0.0.4)
endif
	@echo "==> Triggering release workflow for $(VERSION)..."
	@gh workflow run release.yml -f version=$(VERSION)
	@echo "==> Done! Check GitHub Actions for progress."
	@echo "==> Once PR is created, comment '/release' to publish."
