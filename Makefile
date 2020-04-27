$(GOPATH)/bin/sql-migrate:
	go get -u github.com/rubenv/sql-migrate/...

migrate: $(GOPATH)/bin/sql-migrate
	sql-migrate up -config=_db/config.yml
