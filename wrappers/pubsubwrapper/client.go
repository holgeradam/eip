package pubsubwrapper

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func GetClient(projectID string) (*pubsub.Client, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	return client, nil
}
