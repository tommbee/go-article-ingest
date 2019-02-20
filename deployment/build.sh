#!/bin/bash

echo 'Building...'

echo "Downloading GCP SDK..."
# Downloading gcloud package
curl https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz > /tmp/google-cloud-sdk.tar.gz

# Installing the package
mkdir -p /usr/local/gcloud \
  && tar -C /usr/local/gcloud -xvf /tmp/google-cloud-sdk.tar.gz \
  && /usr/local/gcloud/google-cloud-sdk/install.sh

# Adding the package path to local
ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin

## authenticate with GKE
echo "Authenticating with GKE..."
apt-get install -qq -y gettext
echo $GCLOUD_SERVICE_KEY > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file=${HOME}/gcloud-service-key.json
gsutil cp gs://article-app-storage/${CONFIG_FILENAME} ./${CONFIG_FILENAME}

echo "Building Docker image..."
echo "$DOCKER_PASS" | docker login -u $DOCKER_USER --password-stdin
docker build --build-arg CONFIG_FILENAME=${CONFIG_FILENAME} -t $DOCKER_IMAGE_URL:$CIRCLE_SHA1 .

echo "Pushing to registry..."
docker push $DOCKER_IMAGE_URL:$CIRCLE_SHA1
