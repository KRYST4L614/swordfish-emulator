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
	docker run -d \
		--name gitlab-runner \
		--restart always \
		-v /${PWD}:/${PWD} \
		-v //var/run/docker.sock://var/run/docker.sock \
		gitlab/gitlab-runner:latest
	docker exec -it -w /${PWD} gitlab-runner git config --global --add safe.directory "*"
	docker exec -it -w /${PWD} gitlab-runner gitlab-runner exec docker lint_code
	docker exec -it -w /${PWD} gitlab-runner gitlab-runner exec docker unit_tests
	docker exec -it -w /${PWD} gitlab-runner gitlab-runner exec docker build
