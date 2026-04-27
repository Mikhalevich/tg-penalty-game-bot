#! /usr/bin/env bash

docker build -t bot:$1 -f ./script/docker/bot.Dockerfile .
docker tag bot:$1 registry.digitalocean.com/penalty-bot-registry/bot:$1
docker push registry.digitalocean.com/penalty-bot-registry/bot:$1

docker build -t sqlmigrate:$1 -f ./script/docker/sqlmigrate.Dockerfile .
docker tag sqlmigrate:$1 registry.digitalocean.com/penalty-bot-registry/sqlmigrate:$1
docker push registry.digitalocean.com/penalty-bot-registry/sqlmigrate:$1

docker build -t outboxpoller:$1 -f ./script/docker/outboxpoller.Dockerfile .
docker tag outboxpoller:$1 registry.digitalocean.com/penalty-bot-registry/outboxpoller:$1
docker push registry.digitalocean.com/penalty-bot-registry/outboxpoller:$1

docker build -t gamepoller:$1 -f ./script/docker/gamepoller.Dockerfile .
docker tag gamepoller:$1 registry.digitalocean.com/penalty-bot-registry/gamepoller:$1
docker push registry.digitalocean.com/penalty-bot-registry/gamepoller:$1
