#!/bin/bash

## install helm
echo "Check Helm is installed..."
if [[ $((helm) 2>&1 | grep "command not found" ) ]]; then
    echo "Installing Helm"
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
    chmod 700 get_helm.sh
    ./get_helm.sh
    helm init --client-only
fi
