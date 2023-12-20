.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 && \
	golangci-lint run

.PHONY: test
test:
ifdef CI
	go install gotest.tools/gotestsum@v1.8.2 && \
	gotestsum --junitfile test-report.xml -- -v ./...
else
	go test -v ./...
endif
