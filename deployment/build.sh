#!/bin/bash

echo 'Building...'

## authenticate with GKE
echo "Authenticating with GKE..."
apt-get install -qq -y gettext
echo $GCLOUD_SERVICE_KEY > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file=${HOME}/gcloud-service-key.json
gsutil cp gs://article-app-storage/${CONFIG_FILENAME} ./${CONFIG_FILENAME}

## Change docker registry

echo "Building Docker image..."
echo "$DOCKER_PASS" | docker login -u $DOCKER_USER --password-stdin
docker build --build-arg CONFIG_FILENAME=${CONFIG_FILENAME} -t $DOCKER_IMAGE_URL:$CIRCLE_SHA1 .

echo "Saving docker image to cache..."
mkdir -p docker-cache
docker save -o docker-cache/built-image.tar $DOCKER_IMAGE_URL:$CIRCLE_SHA1
