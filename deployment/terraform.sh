#!/bin/bash

echo $GCLOUD_SERVICE_KEY > ${HOME}/terraform/${GOOGLE_AUTH_FILE}

terraform apply -var "config_file=${GOOGLE_AUTH_FILE}"
