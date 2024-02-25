.PHONY: test-unit
test-unit: ## Run unit tests
	go test -cover ./... -tags=unit  -covermode=count -coverprofile=coverage.out

.PHONY: cover
cover: ## Run unit tests
	go tool cover -func coverage.out

.PHONY: lint
lint: ## Lint the files
	./scripts/lint.sh;

.PHONY: gogen
gogen: ## Regenerate all
	go generate ./...

.PHONY: dep
dep: ## Get the dependencies
	go get -v -d ./...
