ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}

ARG GENERATOR_VERSION=2.3.0

COPY ./oapi-codegen-config.yml /oapi-codegen-config.yml

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v$GENERATOR_VERSION

ENTRYPOINT [ "oapi-codegen", "--config", "/oapi-codegen-config.yml" ]