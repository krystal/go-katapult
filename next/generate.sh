#!/bin/bash

echo "Generating Core and Public clients using OpenAPI specs"

go_version="$(grep goVersion generator-config.yml | cut -d':' -f2 | tr -d '[:space:]')"
generator_version="$(grep generatorVersion generator-config.yml | cut -d':' -f2 | tr -d '[:space:]')"
core_api_version="$(jq -r '.info."x-katapult-version"' katapult-core-openapi.json)"
public_api_version="$(jq -r '.info."x-katapult-version"' katapult-public-openapi.json)"

echo " -> Using Go version: $go_version"
echo " -> Using OpenAPI generator version: $generator_version"

echo " -> Generating Core client (Katapult version: ${core_api_version})..."
docker run \
  --user "$(id -u):$(id -g)" \
  -v "$(pwd):/local" \
  --rm \
  "$(docker build -f Dockerfile --build-arg="GO_VERSION=$go_version" --build-arg="GENERATOR_VERSION=$generator_version" -q .)" \
  -generate types,client \
  -package core \
  -templates /local/templates \
  /local/katapult-core-openapi.json > "./core/core.go"

echo " -> Generating Public client (Katapult version: ${public_api_version})..."
docker run \
  --user "$(id -u):$(id -g)" \
  -v "$(pwd):/local" \
  --rm \
  "$(docker build -f Dockerfile --build-arg="GO_VERSION=$go_version" --build-arg="GENERATOR_VERSION=$generator_version" -q .)" \
  -generate types,client \
  -package public \
  -templates /local/templates \
  /local/katapult-public-openapi.json > "./public/public.go"
