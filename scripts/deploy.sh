#!/bin/sh

usage() {
  cat <<EOF
Description:
    Cloud Functionsに指定したサービスと連携する関数をdeployします。

Usage:
    $(basename ${0}) [service]

Services:
    slack
    dialogflow

Options:
    --help, -h        print this
EOF
}

case $1 in
slack)
  FUNCTION="jimiko-slack"
  ENDPOINT="Slack"
  ;;
dialogflow)
  FUNCTION="jimiko-dialogflow"
  ENDPOINT="Dialogflow"
  ;;
--help | -h)
  usage
  exit 0
  ;;
*)
  echo "slack か dialogflow を指定してください"
  usage
  exit 1
  ;;
esac

if [ -z $FUNCTION ]; then
  echo "function 名が指定されていません"
  exit 1
fi

if [ -z $ENDPOINT ]; then
  echo "endpoint が指定されていません"
  exit 1
fi

gcloud functions deploy $FUNCTION --entry-point $ENDPOINT --trigger-http --runtime=go111 --region=asia-northeast1 --env-vars-file .env.yaml
