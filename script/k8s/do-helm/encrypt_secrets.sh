#! /usr/bin/env bash

sops --encrypt -i script/k8s/do-helm/bot/secrets.yaml
sops --encrypt -i script/k8s/do-helm/redirectionbot/secrets.yaml
sops --encrypt -i script/k8s/do-helm/gamepoller/secrets.yaml
sops --encrypt -i script/k8s/do-helm/outboxpoller/secrets.yaml
