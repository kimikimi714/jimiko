## Install dependencies
.PHONY: deps
deps:
	@go get -v -d

## Run tests
.PHONY: test
test:
	@go test ./...

## Lint
.PHONY: lint
lint:
	@revive ./...

## deploy
.PHONY: deploy
deploy:
	@gcloud functions deploy jimiko-slack-2nd-gen --entry-point Slack \
		--gen2 --trigger-http --region=asia-northeast1 \
		--env-vars-file .env.yaml \
		--runtime=go120 \
		--set-secrets 'SLACK_SIGINING_SECRET=jimiko-slack-signing:latest'
