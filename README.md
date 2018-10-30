# Slack Bot event queue in AWS

This terraform module creates a set of AWS services that allows you to
send Slack event API messages to an AWS SQS queue which can be drained
from within a private network.

The API will answer challenge requests.

### Usage

```
module "swarmbot" {
  source       = "github.com/seanhoughton/slackbot-queue-terraform"
  service_name = "mybot"
}
```

### Inputs

You must define an AWS provider for this module to work.

| Variable     | Description                                                   |
| ------------ | ------------------------------------------------------------- |
| service_name | This will be the prefix on all AWS resources that get created |

### Outputs

| Output    | Description                                        |
| --------- | -------------------------------------------------- |
| event_url | The URL you provide to your Slack bot's events API |