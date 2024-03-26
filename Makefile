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

path:=
.PHONY: gitlab-ci
gitlab-ci: ## Run gitlab CI/CD locally
	./scripts/gitlab-ci.sh;

.PHONY: runner
runner: ## Get the dependencies
	docker run -d --name gitlab-runner --restart always \
		-v ./runner:/etc/gitlab-runner \
		-v /var/run/docker.sock:/var/run/docker.sock \
		gitlab/gitlab-runner:latest
