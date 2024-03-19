#!/bin/bash


go_version=$(grep goVersion generator-config.yml | cut -d':' -f2 | tr -d '[:space:]')
generator_version=$(grep generatorVersion generator-config.yml | cut -d':' -f2 | tr -d '[:space:]')

docker run \
  --user "$(id -u):$(id -g)" \
  -v "$(pwd):/local" \
  --rm \
  "$(docker build -f Dockerfile --build-arg="GO_VERSION=$go_version" --build-arg="GENERATOR_VERSION=$generator_version" -q .)" \
  -generate types,client \
  -package next \
  -templates /local/templates \
   /local/katapult-openapi.json > "./client.go" 
