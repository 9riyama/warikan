GOCMD=go
GOTEST=$(GOCMD) test
GOLANGCI_LINT_VERSION = 1.23.8


$(GOPATH)/bin/sql-migrate:
	go get -u github.com/rubenv/sql-migrate/...

test:
	$(GOTEST) -v ./...

migrate: $(GOPATH)/bin/sql-migrate
	sql-migrate up -config=_db/config.yml

lint:
	golangci-lint run

deps:
	go mod download
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v$(GOLANGCI_LINT_VERSION)
