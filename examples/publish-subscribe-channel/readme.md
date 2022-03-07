# Publish-Subscribe Channel with GCP Pub/Sub

This is a demo implementation of the enterprise integration pattern
[Publish-Subscribe Channel](https://www.enterpriseintegrationpatterns.com/patterns/messaging/PublishSubscribeChannel.html).
It uses Google Cloud Platform's Pub/Sub as the message system.

## How to run

You need to run this as an authenticated GCP user with access to Pub/Sub. Please refer to the
[official documentation](https://cloud.google.com/sdk/gcloud/reference/auth/login) on how to do this.

First you start one or more subscribers from this directoy by running:

```go run ./subscribe -projectID=[PROJECT_ID]```

Then you start the publisher.

```go run ./publish -projectID=[PROJECT_ID]```

The publisher will generate 10 messages out of a sample JSON file. They will be sent to the topic `eip-pub-sub-demo`.
Between each message there's a delay of 3 seconds. You can modify these settings by using the options below.

The subscribers will generate a random UUID for the subscription. By doing this all subscribers will receive every messages.
If you want to split messages between subscribers, see below how to specify the subscription ID and set one subscription
ID for multiple subscribers.

## Command Line Options

```-projectID=[STRING]```

Specify the GCP project ID to be used.

```-topicID=[STRING]```

Specify the topic ID. Defaults to `eip-pub-sub-demo`.

```-howMany=[INT]```

Specify the number of messages to publish. Defaults to 10. (publisher only)

```-delayBy=[DURATION]```

Set the delay between published messages. Defaults to 3 seconds.  (publisher only)

```-subscriptionID=[STRING]```

Set the subscription ID. Defaults to `sub-` and a random UUID.
[GCP limitations](https://cloud.google.com/pubsub/docs/admin#resource_names) apply. (subscriber only)

```-howLong=[DURATION]```

Specify the duration of the subscription being active. Defaults to 30 seconds. (subscriber only)
