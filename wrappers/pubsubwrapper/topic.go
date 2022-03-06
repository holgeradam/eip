package pubsubwrapper

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

func GetTopics(client pubsub.Client) ([]string, error) {
	var topics []string
	ctx := context.Background()

	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to iterate topics: %v", err)
		}
		topics = append(topics, topic.ID())
	}

	return topics, nil
}

func CreateTopic(client pubsub.Client, topicID string) error {
	ctx := context.Background()
	_, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		return fmt.Errorf("Failed to create topic %s: %v", topicID, err)
	}

	return nil
}

func DeleteTopic(client pubsub.Client, topicID string) error {
	ctx := context.Background()
	topic := client.Topic(topicID)
	err := topic.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete topic %s: %v", topicID, err)
	}

	return nil
}
