set -e

readonly service="$1"

echo "start upgrading packages in $service"

if [ "$service" = "pkg" ]; then
    cd "./internal/pkg" && go get -u -t -d -v ./... && go mod tidy
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service" && go get -u -t -d -v ./... && go mod tidy
fi
