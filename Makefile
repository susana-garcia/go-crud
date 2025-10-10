# GOBIN indicates the path where go tooling is put.
GOBIN=$(shell go env GOPATH)/bin

##@ Tools

GOLANGCI_LINT_VER=v2.5.0

golangci-lint: $(GOBIN)/golangci-lint-$(GOLANGCI_LINT_VER) ## installs golangci-lint if not yet installed
	ln -sf $(GOBIN)/golangci-lint-$(GOLANGCI_LINT_VER) $(GOBIN)/golangci-lint
	golangci-lint --version
$(GOBIN)/golangci-lint-$(GOLANGCI_LINT_VER):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VER)
	mv $(GOBIN)/golangci-lint $(GOBIN)/golangci-lint-$(GOLANGCI_LINT_VER)

##@ Development

.PHONY: run
run:
	go run .

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: fmt vet golangci-lint ## Combined target to run linters
	golangci-lint run --timeout 2m0s ./...

.PHONY: start-postgres
start-postgres: ## Starts a local postgres instance
	docker run -d --rm --name postgres -e POSTGRES_USER=gocrud -e POSTGRES_PASSWORD=gocrud -v gocrud-data:/var/lib/postgresql/data -p 5432:5432 postgres:18
	PG_PASSWORD=gocrud docker exec postgres psql -d postgres -h localhost -U gocrud -w -c 'CREATE DATABASE gocrud' || echo "DB already exists - error can be ignored"

.PHONY: stop-postgres
stop-postgres: ## Stops local postgres instance
	@docker container stop $(shell docker container ls -q --filter name=postgres)
