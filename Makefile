EXECUTABLE := api
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
MIGRATIONDIR := store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)

.PHONY: default
default: ${EXECUTABLE}

gobindata: ${MIGRATIONS}
	# Building bindata
	go-bindata -o ${MIGRATIONDIR}/bindata.go -prefix ${MIGRATIONDIR} -pkg migrations ${MIGRATIONDIR}/*.sql

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: gobindata
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

.PHONY: test
test:
	go test ./...

.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI
