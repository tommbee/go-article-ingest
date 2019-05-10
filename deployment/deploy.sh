#!/bin/bash

## install helm
echo "Check Helm is installed..."
if [[ $((helm) 2>&1 | grep "command not found" ) ]]; then
    echo "Installing Helm"
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
    chmod 700 get_helm.sh
    ./get_helm.sh
    helm init --upgrade --kubeconfig ./site-config/kubeconfig
    helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/
fi

gcloud auth activate-service-account --key-file=auth.json

## create namespace
# echo "Creating app namespace..."
# kubectl apply -f ./article-ingest-k8s/namespace.yml

## deploy app
echo "Deploying app via helm..."
helm upgrade -i article-ingest ./article-ingest-k8s \
    --set image.tag=${CIRCLE_SHA1} \
    --set image.repository=${DOCKER_IMAGE_URL} \
    --set sources=${SOURCES} \
    --set server=${SERVER} \
    --set db=${DB} \
    --set configFileLocation=${CONFIG_LOCATION} \
    --set articleCollection=${ARTICLE_COLLECTION} \
    --set dbUser=${DB_USER} \
    --set dbPassword=${DB_PASSWORD} \
    --set authDb=${AUTH_DB} \
    --set dbSsl=${DB_SSL} \
    --kubeconfig ./site-config/kubeconfig
