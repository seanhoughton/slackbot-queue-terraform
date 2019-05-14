IMAGE=seanhoughton/slack-relay
TAG?=test

build:
	go build github.com/seanhoughton/slackbot-queue-terraform/app/cmd/slack-relay

image:
	docker build -t $(IMAGE):$(TAG) .