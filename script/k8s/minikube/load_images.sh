#! /usr/bin/env bash

minikube image build -t bot:$1 -f ./script/docker/bot.Dockerfile .
minikube image build -t sqlmigrate:$1 -f ./script/docker/sqlmigrate.Dockerfile .
minikube image build -t outboxpoller:$1 -f ./script/docker/outboxpoller.Dockerfile .
minikube image build -t gamepoller:$1 -f ./script/docker/gamepoller.Dockerfile .