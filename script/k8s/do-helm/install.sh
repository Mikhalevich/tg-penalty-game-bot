#! /usr/bin/env bash

helm upgrade --install jaeger ./script/k8s/do-helm/jaeger
helm secrets upgrade --install bot ./script/k8s/do-helm/bot -f ./script/k8s/do-helm/bot/secrets.yaml
helm secrets upgrade --install gamepoller ./script/k8s/do-helm/gamepoller -f ./script/k8s/do-helm/gamepoller/secrets.yaml
helm secrets upgrade --install outboxpoller ./script/k8s/do-helm/outboxpoller -f ./script/k8s/do-helm/outboxpoller/secrets.yaml