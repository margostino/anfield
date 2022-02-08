#!make
SHELL = /bin/sh
.DEFAULT: help

-include .env .env.local .env.*.local

# Defaults
BUILD_VERSION ?= SNAPSHOT
IMAGE_NAME := ${DOCKER_REPO}/${SERVICE_NAME}:${BUILD_VERSION}
IMAGE_NAME_LATEST := ${DOCKER_REPO}/${SERVICE_NAME}:latest
DOCKER_COMPOSE = USERID=$(shell id -u):$(shell id -g) docker-compose ${compose-files}
ALL_ENVS := local ci
env ?= local
suse-instances ?= 1
docker-snapshot ?= true

ifndef SERVICE_NAME
$(error SERVICE_NAME is not set)
endif

ifeq (${env}, ci)
compose-files=-f docker-compose.yml -f docker-compose.ci.yml
endif

.PHONY: help
help:
	@echo "Anfield build pipeline"
	@echo ""
	@echo "Usage:"
	@echo "  build                          - Build artifact"
	@echo "  test.unit                      - Run unit tests"
	@echo "  test.integration               - Run integration tests"
	@echo "  test.api                       - Run api tests"
	@echo "  test.resiliency                - Run resiliency tests"
	@echo "  docker.publish                 - Publish docker image (used for internal/external testing purposes) to artifactory. Receives parameter docker-snapshot (default true)"
	@echo "  docker.wait                    - Waits until all docker containers have exited successfully and/or are healthy. Timeout: 180 seconds"
	@echo "  docker.logs                    - Generate one log file per each service running in docker-compose"
	@echo "  git.tag                        - Creates a new tag and pushes it to the git repository. Used to tag the current commit as a released artifact"
	@echo ""
	@echo "  ** The following tasks receive an env parameter to determine the environment they are being executed in. Default env=${env}, possible env values: ${ALL_ENVS}:"
	@echo "  docker.run.dependencies        - Run only SUSE dependencies with docker-compose (default env=${env})". Note that `build` might need to be executed prior.
	@echo "  docker.stop                    - Stop and remove all running containers from this project using docker-compose down (default env=${env})"
	@echo ""
	@echo "Project-level environment variables are set in .env file:"
	@echo "  SERVICE_NAME=anfield"
	@echo "  DOCKER_PROJECT_NAME=anfield"
	@echo "  COMPOSE_PROJECT_NAME=anfield"
	@echo "  DOCKER_REPO="
	@echo "  COMPOSE_HTTP_TIMEOUT=360"
	@echo ""
	@echo "Note: Store protected environment variables in .env.local or .env.*.local"
	@echo ""

.PHONY: build
#TODO

.PHONY: docker.build
#TODO

b.clean: clean
#TODO

.PHONY: test.unit
test.unit:
#TODO

.PHONY: test.integration
test.integration:
#TODO

.PHONY: test.api
test.api:
#TODO

.PHONY: test.resiliency
test.resiliency:
#TODO

.PHONY: docker.run.all
#TODO

.PHONY: docker.run.dependencies
docker.run.dependencies: d.compose.down
	make d.compose.up
	make docker.wait
	docker-compose ps
	docker exec mongo1 /scripts/rs-init.sh

.PHONY: docker.stop
docker.stop: d.compose.down

.PHONY: d.compose.up
d.compose.up:
	$(call DOCKER_COMPOSE) up -d --remove-orphans --build

.PHONY: docker.wait
docker.wait:
	./bin/docker-wait

.PHONY: d.compose.down
d.compose.down:
	$(call DOCKER_COMPOSE) down -v || true
	$(call DOCKER_COMPOSE) rm --force || true
	docker rm "$(docker ps -a -q)" -f || true