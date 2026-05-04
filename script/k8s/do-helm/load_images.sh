#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "load image with version ${VERSION}"

docker build -t bot:${VERSION} -f ./script/docker/bot.Dockerfile .
docker tag bot:${VERSION} registry.digitalocean.com/penalty-bot-registry/bot:${VERSION}
docker push registry.digitalocean.com/penalty-bot-registry/bot:${VERSION}

docker build -t sqlmigrate:${VERSION} -f ./script/docker/sqlmigrate.Dockerfile .
docker tag sqlmigrate:${VERSION} registry.digitalocean.com/penalty-bot-registry/sqlmigrate:${VERSION}
docker push registry.digitalocean.com/penalty-bot-registry/sqlmigrate:${VERSION}

docker build -t outboxpoller:${VERSION} -f ./script/docker/outboxpoller.Dockerfile .
docker tag outboxpoller:${VERSION} registry.digitalocean.com/penalty-bot-registry/outboxpoller:${VERSION}
docker push registry.digitalocean.com/penalty-bot-registry/outboxpoller:${VERSION}

docker build -t gamepoller:${VERSION} -f ./script/docker/gamepoller.Dockerfile .
docker tag gamepoller:${VERSION} registry.digitalocean.com/penalty-bot-registry/gamepoller:${VERSION}
docker push registry.digitalocean.com/penalty-bot-registry/gamepoller:${VERSION}
