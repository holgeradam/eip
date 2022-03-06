package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
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
	var deadLetterTopicID string
	var subID string
	var howLong time.Duration
	var errorQuote float64

	flag.StringVar(&projectID, "projectID", "", "Specify the GCP project ID.")
	flag.StringVar(&topicID, "topicID", "eip-dlc-demo", "Specify a topic ID. Defaults to 'eip-dlc-demo'.")
	flag.StringVar(&deadLetterTopicID, "deadLetterTopicID", "eip-dlc-demo-dead-letters", "Specify a dead lettertopic ID. Defaults to 'eip-dlc-demo-dead-letters'.")
	flag.StringVar(&subID, "subscriptionID", "sub-"+uuid.NewString(), "Specify a subscription ID. Defaults to 'sub-' and a random UUID.")
	flag.DurationVar(&howLong, "howLong", 30*time.Second, "Specify the duration for the active subscription. Defaults to 30s.")
	flag.Float64Var(&errorQuote, "errorQuote", 1.0, "Set the percentage of messages that will not be acked to move to the dead letter queue. Defaults to 1.0.")
	flag.Parse()

	fmt.Println("EIP Pattern Demo: Dead Letter Channel")
	fmt.Println("-------------------------------------")
	fmt.Println("- Subscriber -")
	fmt.Printf("Project ID: %s\n", projectID)
	fmt.Printf("Topic ID: %s\n", topicID)
	fmt.Printf("Dead Letter Topic ID: %s\n", deadLetterTopicID)
	fmt.Printf("Subscription ID: %s\n", subID)
	fmt.Printf("How long: %v\n", howLong)
	fmt.Printf("Error quote: %v\n", errorQuote)
	fmt.Println()

	return subscribe(projectID, topicID, deadLetterTopicID, subID, howLong, errorQuote)
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

func subscribe(projectID string, topicID string, deadLetterTopicID string, subID string, howLong time.Duration, errorQuote float64) error {
	client, err := psw.GetClient(projectID)
	if err != nil {
		return fmt.Errorf("Error creating client: %v", err)
	}

	err = ensureTopicExists(client, topicID)
	if err != nil {
		return err
	}

	fmt.Printf("Creating subscription %s for topic %s lasting %v...\n", subID, topicID, howLong)

	ctx := context.Background()
	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		ExpirationPolicy: 24 * time.Hour,
		Topic:            client.Topic(topicID),
		AckDeadline:      10 * time.Second,
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     "projects/" + projectID + "/topics/" + deadLetterTopicID,
			MaxDeliveryAttempts: 5,
		},
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
			if msg.DeliveryAttempt != nil {
				fmt.Printf("Delivery attempts: %d\n", *msg.DeliveryAttempt)
			}
			f := rand.Float64()
			fmt.Printf("Dice rolled %f.\n", f)
			if errorQuote < f {
				msg.Ack()
				fmt.Println("Message acked.")
			} else {
				fmt.Println("Message not acked.")
			}
			fmt.Println()
		}
	}()

	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msgsIn <- msg
	})

	fmt.Printf("\nReceived %d messages.\n", received)

	return err
}
