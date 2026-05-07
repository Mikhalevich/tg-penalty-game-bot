SHELL = /bin/bash

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin

GOPRIVATE = GOPRIVATE=github.com/Mikhalevich/

LINTER_NAME := golangci-lint
LINTER_VERSION := v2.11.4

APP_TAG := 0.2.0

.PHONY: all build test compose-up compose-down load-test-data vendor install-linter \
lint fmt tools-update generate \
minikube-load-images minikube-apply minikube-delete \
minikube-helm-install minikube-helm-uninstall \
do-load-images do-apply do-delete \
do-helm-load-images do-helm-encrypt-secrets do-helm-install do-helm-uninstall \
install-helm-secrets generate-helm-secrets \

all: build

build:
	go build -mod=vendor -o $(BIN_PATH)/bot ./cmd/bot/main.go
	go build -mod=vendor -o $(BIN_PATH)/gamepoller ./cmd/gamepoller/main.go
	go build -mod=vendor -o $(BIN_PATH)/outboxpoller ./cmd/outboxpoller/main.go
	go build -mod=vendor -o $(BIN_PATH)/botwebhook ./cmd/botwebhook/main.go

test:
	go test ./...

compose-up:
	docker compose -f ./script/docker/docker-compose.yml up --build

compose-down:
	docker compose -f ./script/docker/docker-compose.yml down

load-test-data:
	docker run -it --rm --network host \
		-v ./script/db/dataset/test_data.sql:/script/test_data.sql \
		alpine/psql:17.7 \
		"postgresql://bot:bot@localhost:5432/bot" -f /script/test_data.sql

vendor:
	$(GOPRIVATE) go mod tidy
	$(GOPRIVATE) go mod vendor

install-linter:
	if [ ! -f $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) ]; then \
		echo INSTALLING $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) $(LINTER_VERSION) ; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN)/$(LINTER_VERSION) $(LINTER_VERSION) ; \
		echo DONE ; \
	fi

lint: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) run --config .golangci.yml

fmt: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) fmt --config .golangci.yml

tools-update:
	go get tool

generate:
	$(ENV_PATH) go generate ./...

minikube-load-images:
	./script/k8s/minikube/load_images.sh ${APP_TAG}

minikube-apply: minikube-load-images
	kubectl apply -f ./script/k8s/minikube

minikube-delete:
	kubectl delete -f ./script/k8s/minikube

do-load-images:
	./script/k8s/do/load_images.sh ${APP_TAG}

do-apply:
	kubectl apply -f ./script/k8s/do

do-delete:
	kubectl delete -f ./script/k8s/do

minikube-helm-install:
	./script/k8s/minikube-helm/install.sh ${APP_TAG}

minikube-helm-uninstall:
	./script/k8s/minikube-helm/uninstall.sh

do-helm-load-images:
	./script/k8s/do-helm/load_images.sh ${APP_TAG}

do-helm-encrypt-secrets:
	./script/k8s/do-helm/encrypt_secrets.sh

do-helm-install:
	./script/k8s/do-helm/install.sh ${APP_TAG}

do-helm-uninstall:
	./script/k8s/do-helm/uninstall.sh

install-helm-secrets:
	helm plugin install https://github.com/jkroepke/helm-secrets/releases/download/v4.7.4/secrets-4.7.4.tgz --verify=false

generate-helm-secrets:
	mkdir -p ~/.config/sops/age/
	age-keygen -o ~/.config/sops/age/keys.txt
