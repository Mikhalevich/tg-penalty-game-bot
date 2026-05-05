#! /usr/bin/env bash

helm uninstall --kube-context minikube outboxpoller
helm uninstall --kube-context minikube gamepoller
helm uninstall --kube-context minikube bot
helm uninstall --kube-context minikube jaeger
helm uninstall --kube-context minikube postgres