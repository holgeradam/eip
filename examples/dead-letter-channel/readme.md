# Dead Letter Channel with GCP Pub/Sub

This is a demo implementation of the enterprise integration pattern
[Dead Letter Channel](https://www.enterpriseintegrationpatterns.com/patterns/messaging/DeadLetterChannel.html).
It uses Google Cloud Platform's Pub/Sub as the message system.

## How to run

You need to run this as an authenticated GCP user with access to Pub/Sub. Please refer to the
[official documentation](https://cloud.google.com/sdk/gcloud/reference/auth/login) on how to do this.

First you start the dead letter topic subscriber from this directoy by running:

```go run ./subscribe-dlc -projectID=[PROJECT_ID]```

Second you start the fault simlutating subscriber by running:

```go run ./subscribe -projectID=[PROJECT_ID]```

Then you start the publisher.

```go run ./publish -projectID=[PROJECT_ID]```

The publisher will generate 1 message out of a sample JSON file. It will be sent to the topic `eip-dlc-demo`.

The subscriber will receive the message and by default not ack it. It is configured to allow the minimum of 5 retries.
After that it will forward the message to the dead letter topic `eip-dlc-demo-dead-letters`. The dead letter topic
subscriber will receive it from there.

## Command Line Options

```-projectID=[STRING]```

Specify the GCP project ID to be used.

```-topicID=[STRING]```

Specify the topic ID. Defaults to `eip-pub-sub-demo`.

```-deadLetterTopicID=[STRING]````

Specify the topic ID for the dead letter topic. Defaults to `eip-dlc-demo-dead-letters`. (subscribers only)

```-howMany=[INT]```

Specify the number of messages to publish. Defaults to 10. (publisher only)

```-delayBy=[DURATION]```

Set the delay between published messages. Defaults to 3 seconds.  (publisher only)

```-subscriptionID=[STRING]```

Set the subscription ID. Defaults to `sub-` and a random UUID.
[GCP limitations](https://cloud.google.com/pubsub/docs/admin#resource_names) apply. (subscribers only)

```-howLong=[DURATION]```

Specify the duration of the subscription being active. Defaults to 30 seconds. (subscribers only)
