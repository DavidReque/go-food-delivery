#!/bin/bash

set -e

readonly service="$1"

echo "start building $service"

if [ "$service" = "pkg" ]; then
    cd "./internal/pkg" && go build ./...
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service" && go build ./...
fi
