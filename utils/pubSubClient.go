package utils

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

var pubsubClient *pubsub.Client

func InitPubSubClient() {
	ctx := context.Background()
	var err error
	pubsubClient, err = pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
}

func PublishMessage(topicName string, data []byte) error {
	ctx := context.Background()
	topic := pubsubClient.Topic(topicName)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	// Block until the result is returned and a server-generated ID is returned for the published message
	_, err := result.Get(ctx)
	return err
}
