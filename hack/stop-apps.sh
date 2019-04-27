#!/usr/bin/env bash

DIRECTORY="$(dirname $0)/../"

docker-compose --file "$DIRECTORY/hack/apps.yaml" --project-directory "$DIRECTORY" down
