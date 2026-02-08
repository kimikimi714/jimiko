## Install dependencies
.PHONY: deps
deps:
	@go get -v

## Run tests
.PHONY: test
test:
	@go test ./...

## Lint
.PHONY: lint
lint:
	@revive -set_exit_status ./...

## Modernize Go code
.PHONY: modernize
modernize:
	@modernize -fix ./...

## deploy
.PHONY: deploy
deploy:
	@gcloud functions deploy jimiko-slack-2nd-gen --entry-point Slack \
		--gen2 --trigger-http --region=asia-northeast1 \
		--env-vars-file .env.yaml \
		--runtime=go125 \
		--set-secrets 'SLACK_SIGINING_SECRET=jimiko-slack-signing:latest'
