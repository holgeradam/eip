package router

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

type PubSubMessage struct {
	// Data SampleMessage `json:"data"`
	Data []byte `json:data`
}

type Book struct {
	Author string
	Title  string `json:"title"`
	Isbn13 string `json:"isbn13"`
	Url    string `json:"url"`
}

type SampleMessage struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Book      Book      `json:"book"`
}

type publishRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

var client *pubsub.Client

func init() {
	var err error

	client, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
}

func Route(ctx context.Context, m PubSubMessage) error {
	var sampleMessage SampleMessage
	err := json.Unmarshal([]byte(m.Data), &sampleMessage)
	if err != nil {
		log.Printf("Failed to unmarshal json: %v\n", err)
		return err
	}

	targetTopic := "eip-cbr-demo-read-later"
	if strings.Contains(sampleMessage.Book.Author, "Hohpe") {
		targetTopic = "eip-cbr-demo-read-first"
	}

	log.Printf("Routing message %s to topic %s.\n", sampleMessage.ID, targetTopic)

	jsonMsg, err := json.Marshal(sampleMessage)
	if err != nil {
		log.Printf("Failed to convert to json: %v\n", err)
		return err
	}

	id, err := client.Topic(targetTopic).Publish(ctx, &pubsub.Message{Data: jsonMsg}).Get(ctx)
	if err != nil {
		log.Printf("Failed to send: %v\n", err)
		return err
	}
	log.Printf("Sent message with id %s at %v.\n", id, time.Now().UTC().Format(time.RFC3339Nano))

	return nil
}
