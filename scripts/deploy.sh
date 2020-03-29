#!/bin/sh

while getopts f:p: OPT
do
  case $OPT in
    "f" ) FUNCTION="$OPTARG" ;;
    "p" ) ENDPOINT="$OPTARG" ;;
  esac
done

if [ -z $FUNCTION ] ; then
  echo "function 名が指定されていません"
  exit 1;
fi

if [ -z $ENDPOINT ] ; then
  echo "endpoint が指定されていません"
  exit 1;
fi

gcloud functions deploy $FUNCTION --entry-point $ENDPOINT --trigger-http --runtime=go111 --region=asia-northeast1 --env-vars-file .env.yaml

