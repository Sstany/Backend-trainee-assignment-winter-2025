.PHONY: lint
lint:
	golangci-lint --new-from-rev=master run

.PHONY: coverage
coverage:
	go test -cover ./... 
