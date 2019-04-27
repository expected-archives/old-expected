#!/usr/bin/env bash

DIRECTORY="$(dirname $0)/../"

docker-compose --file "$DIRECTORY/hack/services.yaml" --project-directory "$DIRECTORY" down
