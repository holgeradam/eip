# Content-Based-Router with GCP PubSub and Cloud Functions

This is a demo implementation of the enterprise integration pattern
[Content Based Router](https://www.enterpriseintegrationpatterns.com/patterns/messaging/ContentBasedRouter.html).
It uses Google Cloud Platform's PubSub as the message system and combines it with a Cloud Function to do the content based routing.

## Routing

The example utilizes a Cloud Function as the content based router. It is triggered by the PubSub topic `eip-cbr-demo`. Depending on the content of the `author` property it forwards messages to one of the topics `eip-cbr-demo-read-later` and `eip-cbr-demo-read-first`. I'll let you figure which books go where. ;)

## How to run

You need to run this as an authenticated GCP user with access to PubSub and Cloud Functions. Please refer to the
[official documentation](https://cloud.google.com/sdk/gcloud/reference/auth/login) on how to do this.

To set this up you have to deploy the Cloud Function. Run the following command from the router-function sub-directory:

```gcloud functions deploy Route --runtime=go116 --trigger-topic=eip-cbr-demo --source=./router.go --max-instances=1 --memory=128MB --set-env-vars=[GOOGLE_CLOUD_PROJECT=YOUR_PROJECT_ID] --project=YOUR_PROJECT_ID -region=europe-west1```

Then switch back to the main directory of the example to start the publisher and subscribers.
First you start the subscriber for the topic of books to read later:

```go run ./subscribe -projectID=[PROJECT_ID] -topicID=eip-cbr-demo-read-later```

Second you start the subscriber for the books to read first. It defaults to the corresponding topic:

```go run ./subscribe -projectID=[PROJECT_ID]```

Then you start the publisher.

```go run ./publish -projectID=[PROJECT_ID]```

The publisher will generate 10 message out of a sample JSON file. They will be sent to the topic `eip-cbr-demo`.
The router function will forward the messages and one of the two subscribers will receive them.

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
[GCP limitations](https://cloud.google.com/pubsub/docs/admin#resource_names) apply. (subscribers only)

```-howLong=[DURATION]```

Specify the duration of the subscription being active. Defaults to 30 seconds. (subscribers only)
