.PHONY: start
start: ## Run emulator
	go run cmd/emulator/main.go

.PHONY: unit-test
unit-test: ## Run unit tests
	go test -cover ./... -tags=unit -covermode=count -coverprofile=coverage.out

.PHONY: cover
cover: ## See coverage
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

.PHONY: gitlab-ci
gitlab-ci: ## Run gitlab CI/CD locally
	./scripts/gitlab-ci.sh;

.PHONY: start-runner
start-runner: ## Start runner
	docker-compose -f ./deployments/runner/docker-compose.runner.yaml up -d

.PHONY: stop-runner
stop-runner: ## Stop runner
	docker-compose -f ./deployments/runner/docker-compose.runner.yaml down


.PHONY: oapi-gen
oapi-gen: ## Generate resources structs from OpenAPI spec
	openapi bundle -o docs/openapi/merged.yaml --ext yaml docs/openapi/api.yaml
	oapi-codegen -config=configs/openapi-generator.yaml docs/openapi/merged.yaml
