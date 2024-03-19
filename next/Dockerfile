ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}

ARG GENERATOR_VERSION=2.1.0

RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v$GENERATOR_VERSION

ENTRYPOINT [ "oapi-codegen" ]