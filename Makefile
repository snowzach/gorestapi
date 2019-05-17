EXECUTABLE := api
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
MIGRATIONDIR := store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)
TOOLS := ${GOPATH}/bin/go-bindata \
	${GOPATH}/bin/mockery \
    ${GOPATH}/bin/wire

.PHONY: default
default: ${EXECUTABLE}

${GOPATH}/bin/go-bindata:
	GO111MODULE=off go get -u github.com/go-bindata/go-bindata/...

${GOPATH}/bin/mockery:
	go get github.com/vektra/mockery/cmd/mockery

${GOPATH}/bin/wire:
	go get github.com/google/wire

tools: ${TOOLS}

${MIGRATIONDIR}/bindata.go: ${MIGRATIONS}
	# Building bindata
	go-bindata -o ${MIGRATIONDIR}/bindata.go -prefix ${MIGRATIONDIR} -pkg migrations ${MIGRATIONDIR}/*.sql

cmd/wire_gen.go: cmd/wire.go
	wire ./cmd/...

.PHONY: mocks
mocks: tools
	mockery -dir ./gorestapi -name ThingStore

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools cmd/wire_gen.go ${MIGRATIONDIR}/bindata.go
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

.PHONY: test
test: tools ${MIGRATIONDIR}/bindata.go mocks
	go test -cover ./...

.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI

.PHONY: relocate
relocate:
	@test ${TARGET} || ( echo ">> TARGET is not set. Use: make relocate TARGET=<target>"; exit 1 )
	$(eval ESCAPED_PACKAGENAME := $(shell echo "${PACKAGENAME}" | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_TARGET := $(shell echo "${TARGET}" | sed -e 's/[\/&]/\\&/g'))
	# Renaming package ${PACKAGENAME} to ${TARGET}
	@grep -rlI '${PACKAGENAME}' * | xargs -i@ sed -i 's/${ESCAPED_PACKAGENAME}/${ESCAPED_TARGET}/g' @
	# Complete... 
	# NOTE: This does not update the git config nor will it update any imports of the root directory of this project.
