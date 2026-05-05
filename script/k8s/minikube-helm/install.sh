#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "deploy version ${VERSION}"

helm secrets upgrade --install postgres ./script/k8s/minikube-helm/postgres -f ./script/k8s/minikube-helm/postgres/secrets.yaml
helm upgrade --install jaeger ./script/k8s/minikube-helm/jaeger
helm secrets upgrade --install bot ./script/k8s/minikube-helm/bot -f ./script/k8s/minikube-helm/bot/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
helm secrets upgrade --install gamepoller ./script/k8s/minikube-helm/gamepoller -f ./script/k8s/minikube-helm/gamepoller/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
helm secrets upgrade --install outboxpoller ./script/k8s/minikube-helm/outboxpoller -f ./script/k8s/minikube-helm/outboxpoller/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}