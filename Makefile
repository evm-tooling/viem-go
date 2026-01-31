.PHONY: fmt lint test test-cover verify check

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
