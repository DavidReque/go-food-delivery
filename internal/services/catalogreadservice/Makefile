GOPATH:=$(shell go env GOPATH)

# ==============================================================================
# Go run
.PHONY: run
run:
	go run ./cmd/app/main.go

# ==============================================================================
# Go build
.PHONY: build
build:
	go build ./cmd/main.go

# ==============================================================================
# Go test
.PHONY: test
test:
	go test -cover ./...

# ==============================================================================
# Go lint
.PHONY: lint
lint:
	revive -config revive-config.toml -formatter friendly ./...
	staticcheck ./...
	golangci-lint run ./...

# ==============================================================================
# Go format
.PHONY: format
format:
	golines -m 120 -w --ignore-generated .
	gci write --skip-generated -s standard -s "prefix(github.com/DavidReque/go-food-delivery)" -s default -s blank -s dot --custom-order  .
	gofumpt -l -w .

# ==============================================================================
# Go update
.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod
	go mod tidy

.PHONY: deps-upgrade
deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache

# ==============================================================================
# Docker
.PHONY: docker
docker:
	@docker build -t go-catalogs-read:latest .