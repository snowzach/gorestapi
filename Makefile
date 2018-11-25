EXECUTABLE := api
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
MIGRATIONDIR := store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)
TOOLS := ${GOPATH}/bin/go-bindata \
	${GOPATH}/bin/mockery

.PHONY: default
default: ${EXECUTABLE}

${GOPATH}/bin/go-bindata:
	GO111MODULE=off go get -u github.com/go-bindata/go-bindata/...

${GOPATH}/bin/mockery:
	go get github.com/vektra/mockery/cmd/mockery

tools: ${TOOLS}

${MIGRATIONDIR}/bindata.go: ${MIGRATIONS}
	# Building bindata
	go-bindata -o ${MIGRATIONDIR}/bindata.go -prefix ${MIGRATIONDIR} -pkg migrations ${MIGRATIONDIR}/*.sql

.PHONY: mocks
mocks: tools
	mockery -dir ./gorestapi -name ThingStore

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools ${MIGRATIONDIR}/bindata.go
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

.PHONY: test
test: tools ${MIGRATIONDIR}/bindata.go mocks
	go test -cover ./...

.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI
