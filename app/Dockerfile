FROM golang:alpine AS build-env
COPY . /src
RUN cd /src && go build -mod=vendor github.com/seanhoughton/slackbot-queue-terraform/app/cmd/slack-relay


FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-env /src/slack-relay /slack-relay
EXPOSE 8080
ENV QUEUE VERIFICATION
ENTRYPOINT ["/slack-relay"]
