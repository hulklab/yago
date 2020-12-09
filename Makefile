install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
	go get golang.org/x/tools/cmd/goimports

lint:
	go mod tidy
	goimports -local yago -w .
	gofmt -s -w .
	golangci-lint run -E golint,depguard,gocognit,goconst,gofmt,misspell
