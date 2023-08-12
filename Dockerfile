# Build UI
FROM node:alpine3.18 as build-ui
WORKDIR /build
COPY ui/package.json .
COPY ui/yarn.lock .
RUN yarn
COPY ui/. .
RUN yarn build

# Build API
FROM golang:alpine3.18 AS build

RUN apk add --no-cache --update openssh-client git make
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

SHELL ["/bin/ash", "-c"]

# Handle SSH Private key for github access
# ARG SSH_PRIVATE_KEY
# RUN if [ "${SSH_PRIVATE_KEY}" ]; then \
#   mkdir -p /root/.ssh; \
#   echo "${SSH_PRIVATE_KEY}" > "/root/.ssh/id_rsa"; \
#   chmod 600 ~/.ssh/id_rsa; \
#   ssh-keyscan github.com >> /root/.ssh/known_hosts; \
#   git config --global url."git@github.com:".insteadOf "https://github.com/"; \
# else \
#   echo "This Dockefile requires a build-arg for SSH_PRIVATE_KEY with the contents of an SSH private key to access private repos."; \
#   /bin/false; \
# fi

WORKDIR /build
COPY . .

# Embed the UI
COPY --from=build-ui /build/dist/. embed/public_html

RUN make

FROM alpine:3.18

RUN apk add --no-cache --update tzdata ca-certificates su-exec

# Copy the executable
WORKDIR /app
COPY --from=build /build/build/gorestapi /app/gorestapi
ENTRYPOINT [ "su-exec", "nobody:nobody", "/app/gorestapi" ]
