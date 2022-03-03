package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	var howMany int
	var delayBy time.Duration

	flag.StringVar(&projectID, "projectID", "", "Specify the GCP project ID.")
	flag.StringVar(&topicID, "topicID", "eip-dlc-demo", "Specify a topic ID. Defaults to 'eip-dlc-demo'.")
	flag.IntVar(&howMany, "howMany", 10, "Specify how many messages to send. Defaults to 10.")
	flag.DurationVar(&delayBy, "delayBy", 3*time.Second, "Specify the pause between messages. Defaults to 3s.")
	flag.Parse()

	fmt.Println("EIP Pattern Demo: Dead Letter Channel")
	fmt.Println("-------------------------------------")
	fmt.Println("- Publisher -")
	fmt.Printf("Project ID: %s\n", projectID)
	fmt.Printf("Topic ID: %s\n", topicID)
	fmt.Printf("How many: %d\n", howMany)
	fmt.Printf("Delay by: %v\n", delayBy)
	fmt.Println()

	client, err := psw.GetClient(projectID)
	if err != nil {
		return fmt.Errorf("Error creating client: %v", err)
	}

	return publish(client, topicID, howMany, delayBy)
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
		fmt.Printf("Creating topic: %s\n", topicID)
		err = psw.CreateTopic(*client, topicID)
	}

	return err
}

func publish(client *pubsub.Client, topicID string, deadLetterTopicID string, howMany int, delayBy time.Duration, deleteTopic bool) error {
	err := ensureTopicExists(client, deadLetterTopicID)
	if err != nil {
		return err
	}

	err = ensureTopicExists(client, topicID)
	if err != nil {
		return err
	}

	sampleDataFile, err := ioutil.ReadFile("books.json")
	if err != nil {
		return fmt.Errorf("Error reading sample data file: %v", err)
	}
	var books []Book
	err = json.Unmarshal([]byte(sampleDataFile), &books)
	if err != nil {
		return fmt.Errorf("Error parsing json from sample data file: %v", err)
	}

	fmt.Println("Publishing messages...")
	msgChan := make(chan string)
	resChan := make(chan string)
	errChan := make(chan error)

	go psw.PublishMessage(*client, topicID, msgChan, resChan, errChan)
	go func() {
		for i := 0; i < howMany; i++ {
			message := getRandomMessage(&books)
			jsonMsg, err := json.Marshal(message)
			if err != nil {
				panic(err)
			}
			msgChan <- string(jsonMsg)
			time.Sleep(delayBy)
		}
		close(msgChan)
	}()

	sent, failed := 0, 0
	for res := range resChan {
		sent++
		fmt.Printf("Sent message with id %s at %v.\n", res, time.Now().UTC().Format(time.RFC3339Nano))
	}
	for err := range errChan {
		failed++
		fmt.Printf("Failed to send: %v\n", err)
	}
	fmt.Printf("\nSent %d messages, %d failed.\n", sent, failed)

	if deleteTopic {
		fmt.Printf("Deleting topic: %s\n", topicID)
		err = psw.DeleteTopic(*client, topicID)
		if err != nil {
			return err
		}
	}

	return nil
}

type Book struct {
	Author string
	Title  string `json:"title"`
	Isbn13 string `json:"isbn13"`
	Url    string `json:"url"`
}

type sampleMessage struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Book      Book      `json:"book"`
}

func getRandomMessage(books *[]Book) *sampleMessage {
	randomIndex := rand.Intn(len(*books))
	book := (*books)[randomIndex]

	return &sampleMessage{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Message:   "Looking for a good book to read? Try this one.",
		Book:      book,
	}
}

func getRandomString() string {
	return "abc"
}
