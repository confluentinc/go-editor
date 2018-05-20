deps:
	@which golint 2>/dev/null || go get -u github.com/golang/lint/golint

test:
	@gofmt -w .
	@golint -set_exit_status ./...
	@go vet ./...
	@go test -race -coverprofile=coverage.txt -covermode=atomic
