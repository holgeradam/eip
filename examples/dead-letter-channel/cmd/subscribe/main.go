
// func newSubscription(ctx *context.Context, client *pubsub.Client, topicID string, subID string, howLong time.Duration, handle func(chan *pubsub.Message, chan int)) error {
// 	err := ensureTopicExists(client, topicID)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("Creating subscription %s for topic %s lasting %v...\n", subID, topicID, howLong)

// 	sub, err := client.CreateSubscription(*ctx, subID, pubsub.SubscriptionConfig{
// 		ExpirationPolicy: 24 * time.Hour,
// 		Topic:            client.Topic(topicID),
// 	})
// 	if err != nil {
// 		return fmt.Errorf("Error creating subscription: %v", err)
// 	}
// 	cctx, cancel := context.WithTimeout(*ctx, howLong)
// 	defer cancel()

// 	msgsIn := make(chan *pubsub.Message)
// 	defer close(msgsIn)

// 	count := make(chan int)
// 	go handle(msgsIn, count)
// 	received := <-count

// 	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
// 		msgsIn <- msg
// 	})

// 	fmt.Printf("\nReceived %d messages.\n", received)

// 	return err
// }

// func subscribe(ctx *context.Context, client *pubsub.Client, topicID string, deadLetterTopicID string, subID string, howLong time.Duration, errorQuote float64) error {
// 	err := newSubscription(ctx, client, topicID, subID, howLong, func(msgsIn chan *pubsub.Message, count chan int) {
// 		received := 0
// 		for msg := range msgsIn {
// 			received++
// 			fmt.Printf("Received message %s at %v:\nPublish date: %v\n%q\n\n",
// 				msg.ID,
// 				time.Now().UTC().Format(time.RFC3339Nano),
// 				msg.PublishTime.UTC().Format(time.RFC3339Nano),
// 				string(msg.Data))
// 			if len(msg.Attributes) > 0 {
// 				fmt.Println("Attributes:")
// 				for key, value := range msg.Attributes {
// 					fmt.Printf("%s = %s", key, value)
// 				}
// 			}
// 			if errorQuote < rand.Float64() {
// 				msg.Ack()
// 			} else {
// 				msg.Nack()
// 			}
// 		}
// 		count <- received
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return newSubscription(ctx, client, deadLetterTopicID, subID, howLong, func(msgsIn chan *pubsub.Message, count chan int) {
// 		received := 0
// 		for msg := range msgsIn {
// 			received++
// 			fmt.Printf("Received dead letter message %s at %v:\nPublish date: %v\n%q\n\n",
// 				msg.ID,
// 				time.Now().UTC().Format(time.RFC3339Nano),
// 				msg.PublishTime.UTC().Format(time.RFC3339Nano),
// 				string(msg.Data))
// 			if len(msg.Attributes) > 0 {
// 				fmt.Println("Attributes:")
// 				for key, value := range msg.Attributes {
// 					fmt.Printf("%s = %s", key, value)
// 				}
// 			}
// 			msg.Ack()
// 		}
// 		count <- received
// 	})
// }
