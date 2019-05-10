#!/bin/bash

echo 'Getting config...'

## authenticate with GKE
echo "Authenticating with GKE..."
apt-get install -qq -y gettext
echo $GCLOUD_SERVICE_KEY > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file=${HOME}/gcloud-service-key.json

echo "Saving config to cache..."
mkdir -p site-config
echo $GCLOUD_SERVICE_KEY > ./site-config/auth.json
gsutil cp gs://article-app-storage/${CONFIG_FILENAME} ./site-config/${CONFIG_FILENAME}
gsutil cp gs://article-app-storage/kubeconfig ./site-config/kubeconfig
