#!/bin/bash

ssh -o StrictHostKeyChecking=no root@$HOST_ADDRESS << 'ENDSSH'
    docker pull $DOCKER_IMAGE_URL:$CIRCLE_SHA1
    docker tag $DOCKER_IMAGE_URL:$CIRCLE_SHA1 dokku/$DOKKU_APP_NAME:$CIRCLE_SHA1
    dokku tags:deploy $DOKKU_APP_NAME $CIRCLE_SHA1
ENDSSH
