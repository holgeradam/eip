package pubsubwrapper

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
)

func PublishMessage(client pubsub.Client, topicID string, msgChan <-chan string, resChan chan<- string, errChan chan<- error) {
	ctx := context.Background()
	topic := client.Topic(topicID)

	topic.PublishSettings.CountThreshold = 5
	topic.PublishSettings.DelayThreshold = 100 * time.Millisecond

	for msg := range msgChan {
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(msg),
		})
		id, err := result.Get(ctx)
		if err != nil {
			errChan <- err
		} else {
			resChan <- id
		}
	}
	close(resChan)
	close(errChan)
}
