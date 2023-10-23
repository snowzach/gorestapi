# Build UI
FROM node:alpine3.18 as build-ui
WORKDIR /build
COPY ui/package.json .
COPY ui/yarn.lock .
RUN yarn
COPY ui/. .
RUN yarn build

# Build API
FROM golang:alpine3.18 AS build-api

RUN apk add --no-cache --update openssh-client git make
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

SHELL ["/bin/ash", "-c"]

# Setup SSH for private repos
RUN mkdir -m 0700 ~/.ssh \
    && ssh-keyscan github.com >> ~/.ssh/known_hosts \
    && git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /build
COPY . .

# Embed the UI
COPY --from=build-ui /build/dist/. embed/public_html

RUN --mount=type=ssh make

FROM alpine:3.18

RUN apk add --no-cache --update tzdata ca-certificates su-exec

# Copy the executable
WORKDIR /app
COPY --from=build-api /build/build/gorestapi /app/gorestapi
ENTRYPOINT [ "su-exec", "nobody:nobody", "/app/gorestapi" ]
