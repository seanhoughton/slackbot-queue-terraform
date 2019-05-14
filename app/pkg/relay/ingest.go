package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
)

type ingestEvent struct {
	AppKey string      `json:"app_key"`
	Event  interface{} `json:"event"`
}

// Event is an event from the events API
type Event struct {
	AppKey string
	Event  slackevents.EventsAPIEvent
}

func processMessage(events chan<- *Event, message *sqs.Message, token string) error {
	ingest := ingestEvent{}
	err := json.Unmarshal(json.RawMessage(*message.Body), &ingest)
	if err != nil {
		return fmt.Errorf("Malformed event: %v", err)
	}

	// re-encode it so the parse event can work
	ingestEventStr, err := json.Marshal(ingest.Event)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	event, err := slackevents.ParseEvent(json.RawMessage(ingestEventStr), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: token}))
	if err != nil {
		return fmt.Errorf("Failed to decode message: %v", err)
	}

	events <- &Event{AppKey: ingest.AppKey, Event: event}
	return nil
}

// Poll starts an infinite polling loop listening for messages and sending them to the returned channel
func Poll(ctx context.Context, queueURL string, verificationToken string) <-chan *Event {
	eventChan := make(chan *Event)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	queue := sqs.New(sess)

	go func() {
		for {
			items, err := queue.ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				},
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            &queueURL,
				MaxNumberOfMessages: aws.Int64(10),
				VisibilityTimeout:   aws.Int64(20), // 20 seconds
				WaitTimeSeconds:     aws.Int64(0),
			})

			if err != nil {
				log.Errorf("Failed to receive message: %v", err)
				continue
			}

			if len(items.Messages) == 0 {
				log.Debug("No new messages available")
				time.Sleep(time.Second)
				continue
			}

			log.Infof("Processing %d message(s)...", len(items.Messages))
			for _, message := range items.Messages {
				log.Debugf("Processing %s", message)

				err := processMessage(eventChan, message, verificationToken)
				if err != nil {
					log.Errorf("Failed to process message: %v", err)
				}

				_, err = queue.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &queueURL,
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					log.Errorf("Failed to delete processed message: %v", err)
				}
			}
		}
	}()

	return eventChan
}
