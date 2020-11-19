EXECUTABLE := gorestapicmd
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
EMBEDDIR := embed
TOOLS := ${GOPATH}/bin/go-bindata \
	${GOPATH}/bin/mockery \
	${GOPATH}/bin/swagger
SWAGGERSOURCE = $(wildcard gorestapi/*.go) \
	$(wildcard gorestapi/thingrpc/*.go)

.PHONY: default
default: ${EXECUTABLE}

tools: ${TOOLS}

${GOPATH}/bin/go-bindata:
	GO111MODULE=off go get -u github.com/go-bindata/go-bindata/go-bindata

${GOPATH}/bin/mockery:
	go get github.com/vektra/mockery/cmd/mockery 

${GOPATH}/bin/swagger:
	wget -O ${GOPATH}/bin/swagger https://github.com/go-swagger/go-swagger/releases/download/v0.25.0/swagger_linux_amd64
	chmod 755 ${GOPATH}/bin/swagger

${EMBEDDIR}/bindata.go: tools $(wildcard embed/postgres_migrations/*.sql) $(wildcard embed/public/api-docs/*)
	# Building bindata
	go-bindata -o ${EMBEDDIR}/bindata.go -prefix ${EMBEDDIR} -pkg embed ${EMBED} embed/postgres_migrations/... embed/public/...

.PHONY: swagger
swagger: tools ${SWAGGERSOURCE}
	swagger generate spec --scan-models -o embed/public/api-docs/swagger.json

embed/public/api-docs/swagger.json: tools ${SWAGGERSOURCE}
	swagger generate spec --scan-models -o embed/public/api-docs/swagger.json 

.PHONY: mocks
mocks: tools
	mockery -dir ./gorestapi -name ThingStore

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools ${EMBEDDIR}/bindata.go
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

.PHONY: test
test: tools ${EMBEDDIR}/bindata.go mocks
	go test -cover ./...

.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI

.PHONY: lint
lint:
	docker run --rm -v ${PWD}:/app -e -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v --timeout 5m

.PHONY: hadolint
hadolint:
	docker run -it --rm -v ${PWD}/Dockerfile:/Dockerfile hadolint/hadolint:latest hadolint --ignore DL3018 Dockerfile

.PHONY: relocate
relocate:
	@test ${TARGET} || ( echo ">> TARGET is not set. Use: make relocate TARGET=<target>"; exit 1 )
	$(eval ESCAPED_PACKAGENAME := $(shell echo "${PACKAGENAME}" | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_TARGET := $(shell echo "${TARGET}" | sed -e 's/[\/&]/\\&/g'))
	# Renaming package ${PACKAGENAME} to ${TARGET}
	@grep -rlI '${PACKAGENAME}' * | xargs -i@ sed -i 's/${ESCAPED_PACKAGENAME}/${ESCAPED_TARGET}/g' @
	# Complete... 
	# NOTE: This does not update the git config nor will it update any imports of the root directory of this project.

