ARG GO_VERSION=1.21
FROM golang:$GO_VERSION

ARG GENERATOR_VERSION=v2.0.0

RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@${GENERATOR_VERSION}

ENTRYPOINT [ "oapi-codegen" ]