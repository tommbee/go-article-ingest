#!/bin/bash

export HELM_HOME=/home/circleci/project/.helm

## install helm
echo "Check Helm is installed..."
if [[ $((helm) 2>&1 | grep "command not found" ) ]]; then
    echo "Installing Helm"
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
    chmod 700 get_helm.sh
    ./get_helm.sh
    helm init --client-only
    helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/
fi
