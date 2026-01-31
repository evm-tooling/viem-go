.PHONY: fmt lint test test-cover verify

fmt:
	gofmt -w .
	goimports -w .

lint: fmt
	golangci-lint run

verify:
	golangci-lint config verify

test:
	go test -v -race $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

test-cover:
	go test -v -race -coverprofile=coverage.out $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)