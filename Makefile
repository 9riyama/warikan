GOCMD=go
GOTEST=$(GOCMD) test

$(GOPATH)/bin/sql-migrate:
	go get -u github.com/rubenv/sql-migrate/...

test:
	$(GOTEST) -v ./...

migrate: $(GOPATH)/bin/sql-migrate
	sql-migrate up -config=_db/config.yml
