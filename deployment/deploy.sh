#!/bin/bash

echo 'Deploying...'

export NAMESPACE=article-app
export MONITORING_NAMESPACE=monitoring

## install helm
echo "Get Helm installed..."
if [[ $((helm) 2>&1 | grep "command not found" ) ]]; then
    echo "Installing Helm"
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
    chmod 700 get_helm.sh
    ./get_helm.sh
fi

## authenticate with GKE
echo "Authenticating with GKE..."
apt-get install -qq -y gettext
echo $GCLOUD_SERVICE_KEY > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file=${HOME}/gcloud-service-key.json
gcloud --quiet config set project ${GOOGLE_PROJECT_ID}
gcloud --quiet config set compute/zone ${GOOGLE_COMPUTE_ZONE}
gcloud --quiet container clusters get-credentials ${GOOGLE_CLUSTER_NAME}

## install prometheus
echo "Installing prometheus..."
helm init
helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/
helm install --name prometheus-operator --namespace ${MONITORING_NAMESPACE} stable/prometheus-operator
helm install --name kube-prometheus --namespace ${MONITORING_NAMESPACE} coreos/kube-prometheus --set global.rbacEnable=true

## create namespace
echo "Creating app namespace..."
kubectl apply -f ./article-ingest-k8s/namespace.yml

## deploy helm chart
echo "Deploying helm chart..."
helm upgrade -i article-ingest ./article-ingest-k8s \
                --set image.repository=${DOCKER_IMAGE_URL} \
                --set image.tag=${CIRCLE_SHA1} \
                --set sources="${SOURCES}" \
                --set server=${SERVER} \
                --set db=${DB} \
                --set configFileLocation="${CONFIG_LOCATION}" \
                --set articleCollection=${ARTICLE_COLLECTION} \
                --set dbUser=${DB_USER} \
                --set dbPassword=${DB_PASSWORD} \
                --set authDb=${AUTH_DB} \
                --set dbSsl=${DB_SSL} \
                --namespace ${NAMESPACE}
