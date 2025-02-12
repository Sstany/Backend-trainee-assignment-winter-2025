.PHONY: lint
lint:
	golangci-lint --new-from-rev=main run

.PHONY: coverage
coverage:
	go test -cover ./... 
