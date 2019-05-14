# A Slack Events API Relay Service

This repo contains all the components required to run a Slack Event API relay inside a private network so
bots don't need to expose public-facing endpoints.

The pipeline uses an AWS queue to buffer the events so they can be polled.

```
API Gateway -> Lambda -> SQS <- relay -> bot(s)
```


## Setting up the AWS side with Terraform

The `/terraform` folder contains the required terraform configurations and modules to create the pipeline. You must have an AWS account
and create the `/terraform/terraform.tfvars` file containing your credentials. The `service_name` will be prefixed on every AWS resource.

```
access_key = "xxxxxxxxxxxxxxxxxxxx"
secret_key = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
service_name = "my-slack-relay"
```

## Configuring your bot

First, pick a secret key to use for your app. This can be any string that is valid as a URL path segment. For example: "mybot1234xyz"

Look at the generated AWS API Gateway to find the public deployment address. Each Slack app should also pick a secret key. This key is used for polling and provides security against unauthorized access of relay messages.

Enter the URL with the following format: `{gateway_url}/{my_app_key}/event`. For example:

```
https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/prod/mybot1234xyz/event
```

## Running the relay

Run the relay application from inside your LAN. This becomes the server for listening to events.

| Argument     | Environment Variable | Description                           |
| ------------ | -------------------- | ------------------------------------- |
| addr         |                      | What to listen on in `ip:port` format |
| queue        | `QUEUE`              | The full URL of the AWS event queue   |
| verification | `VERIFICATION`       | The Slack app's verification token    |


```
slack-relay -port :8080 -queue https://sqs.us-east-1.amazonaws.com/xxxxxxxxxxxx/my-slack-relay.fifo -verification xxxxxxxxxxxxxxxxxxxxxxxx
```

## Subscribing to events

The relay exposes a websocket for each app key at the `/{key}/ws` path. For example:

```
wscat -c http://localhost:8080/mybot1234xyz/ws
```

Notes:
* Multiple subscribers to the same app key will get **duplicate** events. This is a demultiplexer not a queue.
* Events that arrive when no subscriptions exist for that app key will be discarded.

