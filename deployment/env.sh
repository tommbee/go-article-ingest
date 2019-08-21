#!/bin/bash

ENVVAR_LIST=''
while read LINE; do
    KEY=''
    VALUE=''

    LINE=`echo $LINE | sed 's/\"//g'`

    KEY=`echo $LINE | sed 's/=/ /' | cut -d ' ' -f1`
    VALUE=`echo $LINE | sed 's/=/ /' | cut -d ' ' -f2`

    if [[ $VALUE == "~" ]]; then
        VALUE=$(echo ${!KEY})
    fi;

    ENVVAR_LIST+="${KEY}=\"${VALUE}\" "
done < $1

DOKKU_COMMAND="config:set ${DOKKU_APP_NAME} ${ENVVAR_LIST}"
DOKKU_COMMAND=`echo $DOKKU_COMMAND | sed "$ s/$//g"`
FINAL_COMMAND="curl -i -X POST -d 'cmd=${DOKKU_COMMAND}' -H 'Api-Key: ${DOKKU_API_KEY}' -H 'Api-Secret: ${DOKKU_API_SECRET}' ${DOKKU_API_ENDPOINT}"

eval $FINAL_COMMAND
