.PHONY: up down check-env

DIR := $(dir $(lastword $(MAKEFILE_LIST)))

up: check-env ## perform forward schema migrations
	migrate -database "mongodb://$(MONGODB_USER):$(MONGODB_PASS)@$(MONGODB_HOST):$(MONGODB_PORT)/$(MONGODB_DATABASE)?w=majority&wtimeoutMS=5000" -path $(MIGRATIONS_DIR) up $(step)

down: check-env ## perform backward schema migrations
	migrate -database "mongodb://$(MONGODB_USER):$(MONGODB_PASS)@$(MONGODB_HOST):$(MONGODB_PORT)/$(MONGODB_DATABASE)?w=majority&wtimeoutMS=5000" -path $(MIGRATIONS_DIR) down $(step)

reset: check-env down up ## reset the database

check-env:
ifndef MONGODB_HOST
  $(error MONGODB_HOST is undefined)
endif
ifndef MONGODB_USER
  $(error MONGODB_USER is undefined)
endif
ifndef MONGODB_PASS
  $(error MONGODB_PASS is undefined)
endif
ifndef MONGODB_PORT
  $(error MONGODB_PORT is undefined)
endif
ifndef MONGODB_DATABASE
  $(error MONGODB_DATABASE is undefined)
endif
ifndef MIGRATIONS_DIR
  $(error MIGRATIONS_DIR is undefined)
endif

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
