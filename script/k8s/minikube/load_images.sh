#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "load image with version ${VERSION}"

minikube image build -t bot:${VERSION} -f ./script/docker/bot.Dockerfile .
minikube image build -t sqlmigrate:${VERSION} -f ./script/docker/sqlmigrate.Dockerfile .
minikube image build -t outboxpoller:${VERSION} -f ./script/docker/outboxpoller.Dockerfile .
minikube image build -t gamepoller:${VERSION} -f ./script/docker/gamepoller.Dockerfile .