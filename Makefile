.PHONY: dependencies
dependencies:
	@go mod download
	@go mod tidy
	@./tools/install.sh

.PHONY: dump-env
dump-env:
	@envsubst < env.yml.tmpl > env.yml

.PHONY: deploy
deploy: deploy-webhook deploy-pubsub

.PHONY: deploy-webhook
deploy-webhook: dependencies dump-env
	@gcloud functions deploy $(GCLOUD_LINE_WEBHOOK_FUNCTION_NAME) \
	--entry-point WebhookHandler \
	--trigger-http \
	--runtime go113 \
	--memory 128MB \
	--source . \
 	--env-vars-file=env.yml

.PHONY: deploy-pubsub
deploy-pubsub: dependencies dump-env
	@gcloud functions deploy $(GCLOUD_DRIVE_UPLOAD_FUNCTION_NAME) \
	--entry-point HandlePubSub \
	--trigger-topic="${GCP_PUBSUB_TOPIC}" \
	--runtime go113 \
	--memory 128MB \
	--source . \
 	--env-vars-file=env.yml
