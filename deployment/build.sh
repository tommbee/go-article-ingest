#!/bin/bash

echo 'Building...'

## Copy contents of Google bucket down to ${CONFIG_FILENAME}
## authenticate with GKE
echo "Authenticating with GKE..."
apt-get install -qq -y gettext
echo $GCLOUD_SERVICE_KEY > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file=${HOME}/gcloud-service-key.json
gsutil cp gs://article-app-storage/${CONFIG_FILENAME} ./${CONFIG_FILENAME}

echo "Building Docker image..."
docker login -u $DOCKER_USER -p $DOCKER_PASS
docker build --build-arg CONFIG_FILENAME=${CONFIG_FILENAME} -t $DOCKER_IMAGE_URL:$CIRCLE_SHA1 .

echo "Pushing to registry..."
docker push $DOCKER_IMAGE_URL:$CIRCLE_SHA1
