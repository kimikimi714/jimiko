#!/bin/sh

usage() {
  cat <<EOF
Description:
    Cloud Functionsに指定したサービスと連携する関数をdeployします。

Usage:
    $(basename ${0}) [service]

Services:
    slack (default)

Options:
    --help, -h        print this
EOF
}

case $1 in
--help | -h)
  usage
  exit 0
  ;;
*)
  FUNCTION="jimiko-slack-2nd-gen"
  ENDPOINT="Slack"
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

gcloud functions deploy $FUNCTION --entry-point $ENDPOINT \
  --gen2 --trigger-http --region=asia-northeast1 \
  --env-vars-file .env.yaml \
  --runtime=go120 --set-secrets 'SLACK_SIGINING_SECRET=jimiko-slack-signing:latest'
