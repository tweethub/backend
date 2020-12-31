.PHONY: check-direnv

LINTER_VERSION := v1.33.0

check: test lint doc

clean: stop-dependencies ## clean docker environment
	docker-compose rm -f nats-streaming mongo redis tweethub.io tweets statistics

doc: ## generate API documentation
	swag init -g doc.go
	redoc-cli bundle -o docs/index.html docs/swagger.json

gen: check-direnv ## generate protobufs and mocks
	go generate ./...

lint: check-direnv ## run lint check
	prototool lint api/events
	golangci-lint run --fix

show-doc: doc ## show the API documentation
	redoc-cli serve docs/swagger.json

setup: check-direnv ## install required development utilities
	go mod tidy
	go get github.com/golang/mock/mockgen
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/uber/prototool/cmd/prototool
	go get github.com/swaggo/swag/cmd/swag
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b bin $(LINTER_VERSION)
	curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && sudo install skaffold /usr/local/bin/; rm skaffold
	yarn install

start: stop-dependencies start-dependencies ## start the backend locally
	docker-compose up --build tweets statistics

start-dependencies: ## start third party services
	docker-compose pull nats-streaming mongo redis tweethub.io
	docker-compose up -d --build nats-streaming mongo redis tweethub.io

start-frontend: ## start the frontend locally
	yarn --cwd $$(dirname `pwd`)/frontend start

stop: stop-dependencies ## stop the backend
	docker-compose stop tweets statistics

stop-dependencies: ## stop third party services
	docker-compose stop nats-streaming mongo redis tweethub.io

test: check-direnv ## run the unit tests
	go test ./... -race -coverprofile .test_coverage.txt

test-coverage: test ## show test coverage
	go tool cover -html=.test_coverage.txt

check-direnv:
ifeq ($(DIRENV_DIR),)
	@echo "direnv or it's shell integration is not installed correctly"
	@exit 1
endif

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
