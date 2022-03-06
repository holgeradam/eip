package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	psw "github.com/holgeradam/eip/wrappers/pubsubwrapper"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var projectID string
	var topicID string
	var subID string
	var howLong time.Duration

	flag.StringVar(&projectID, "projectID", "", "Specify the GCP project ID.")
	flag.StringVar(&topicID, "topicID", "eip-pub-sub-demo", "Specify a topic ID. Defaults to 'eip-pub-sub-demo'.")
	flag.StringVar(&subID, "subscriptionID", "sub-"+uuid.NewString(), "Specify a subscription ID. Defaults to 'sub-' and a random UUID.")
	flag.DurationVar(&howLong, "howLong", 30*time.Second, "Specify the duration for the active subscription. Defaults to 30s.")
	flag.Parse()

	fmt.Println("EIP Pattern Demo: Publish Subscribe Channel")
	fmt.Println("-------------------------------------------")
	fmt.Println("- Subscriber -")
	fmt.Printf("Project ID: %s\n", projectID)
	fmt.Printf("Topic ID: %s\n", topicID)
	fmt.Printf("Subscription ID: %s\n", subID)
	fmt.Printf("How long: %v\n", howLong)
	fmt.Println()

	return subscribe(projectID, topicID, subID, howLong)
}

func ensureTopicExists(client *pubsub.Client, topicID string) error {
	topicExists := false

	topics, err := psw.GetTopics(*client)
	if err != nil {
		return err
	}

	for _, t := range topics {
		if t == topicID {
			topicExists = true
		}
	}

	if topicExists {
		fmt.Printf("Topic %s already exists. Skipping create...\n", topicID)
	} else {
		fmt.Printf("Creating topic %s...\n", topicID)
		err = psw.CreateTopic(*client, topicID)
	}

	return err
}

func subscribe(projectID string, topicID string, subID string, howLong time.Duration) error {
	client, err := psw.GetClient(projectID)
	if err != nil {
		return fmt.Errorf("Error creating client: %v", err)
	}

	err = ensureTopicExists(client, topicID)
	if err != nil {
		return err
	}

	fmt.Printf("Creating subscription %s for topic %s lasting %v...\n", subID, topicID, howLong)
	fmt.Println()

	ctx := context.Background()
	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		ExpirationPolicy: 24 * time.Hour,
		Topic:            client.Topic(topicID),
	})
	sub.ReceiveSettings.MaxExtension = -1 * time.Second

	if err != nil {
		return fmt.Errorf("Error creating subscription: %v", err)
	}
	cctx, cancel := context.WithTimeout(ctx, howLong)
	defer cancel()

	msgsIn := make(chan *pubsub.Message)
	defer close(msgsIn)

	received := 0
	go func() {
		for msg := range msgsIn {
			received++
			fmt.Printf("Received message %s at %v:\nPublish date: %v\n%q\n",
				msg.ID,
				time.Now().UTC().Format(time.RFC3339Nano),
				msg.PublishTime.UTC().Format(time.RFC3339Nano),
				string(msg.Data))
			if len(msg.Attributes) > 0 {
				fmt.Println("Attributes:")
				for key, value := range msg.Attributes {
					fmt.Printf("%s = %s\n", key, value)
				}
			}
			msg.Ack()
			fmt.Println()
		}
	}()

	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msgsIn <- msg
	})

	fmt.Printf("\nReceived %d messages.\n", received)

	return err
}
