LINUXAMD64 = CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# try refresh envvar from .env
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# common commands

envvar-exists-%:
	@if [ -z '${${*}}' ]; then echo 'ERROR: variable $* not set' && exit 1; fi

cmd-exists-%:
	@hash $(*) > /dev/null 2>&1 || \
		(echo "ERROR: '$(*)' must be installed and available on your PATH."; exit 1)

.PHONY: envvar-exists-% cmd-exists-%

# custom commands

deps:
	go mod tidy
	go mod download

build: deps
	$(LINUXAMD64) go build -o unico-challenge .

up: envvar-exists-DATABASE_URL
	@migrate -path ./migrations -database ${DATABASE_URL} -verbose up

up-to-%: envvar-exists-DATABASE_URL
	@migrate -path ./migrations -database ${DATABASE_URL} -verbose up $(*)

down-to-%: envvar-exists-DATABASE_URL
	@migrate -path ./migrations -database ${DATABASE_URL} -verbose down $(*)

serve:
	@go run main.go serve

test:
	@go test -count=1 -cover -race ./...

cover:
	@go test -coverprofile=cover.out.tmp ./...
	@cat cover.out.tmp | grep -v "mock.go" > cover.out
	@go tool cover -html=cover.out -o=cover.html

mockgen:
	@mockgen -source ./repository/feiralivre/feiralivre.go -destination ./repository/feiralivre/mock.go -package feiralivre

lint:
	@golangci-lint run ./...

clean:
	@if test -f "cover.html" ; then rm cover.html ; fi
	@if test -f "cover.out" ; then rm cover.out ; fi
	@if test -f "cover.out.tmp" ; then rm cover.out.tmp ; fi
	@if test -f "logs.txt" ; then rm logs.txt ; fi
	@if test -f "unico-challenge" ; then rm unico-challenge ; fi

.PHONY: deps build up up-to-% down-to-% serve test cover mockgen lint clean
