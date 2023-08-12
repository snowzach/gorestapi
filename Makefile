EXECUTABLE := gorestapi
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
TOOLS := ${GOPATH}/bin/mockery \
	${GOPATH}/bin/swag
SWAGGERSOURCE = $(wildcard gorestapi/*.go) \
	$(wildcard gorestapi/mainrpc/*.go)

.PHONY: default
default: ${EXECUTABLE}

tools: ${TOOLS}

${GOPATH}/bin/mockery:
	go install github.com/vektra/mockery/v2@latest

${GOPATH}/bin/swag:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swagger
swagger: tools ${SWAGGERSOURCE}
	swag init --dir . --parseDependency --parseDepth 1 --generalInfo gorestapi/swagger.go --exclude embed --output embed/public_html/api-docs
	rm embed/public_html/api-docs/docs.go
	rm embed/public_html/api-docs/swagger.yaml
	
.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools swagger
	# Compiling...
	mkdir -p build
	go build -ldflags "-X github.com/snowzach/golib/version.Executable=${EXECUTABLE} -X github.com/snowzach/golib/version.GitVersion=${GITVERSION}" -o build/${EXECUTABLE}

.PHONY: mocks
mocks: tools
	mockery --dir ./gorestapi --name GRStore 

.PHONY: test
test: tools mocks
	go test -cover ./...

.PHONY: lint
lint:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.54.0 golangci-lint run -v --timeout 5m

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

