# lp2gd

<img src="https://user-images.githubusercontent.com/536667/75862456-55e2c980-5e42-11ea-9da8-23f2f2894882.png" width="400px">

## What is lp2gd ?

lp2gd is LINE Photo to Google Drive

1. LINE Webhook to Google Cloud Function (GCF)
1. GCF publish message to Cloud Pub/Sub (Pub/Sub) when web hook message type is image
1. GCF get content image from line and upload Google Drive when subscribe message from Pub/Sub

## How to use

- Create new LINE channel and edit information
- Create GCP project
- Create service account and share Google Drive
- Define env
- make deploy

## Env

| env name                              | description                                       |
| ------------------------------------- | ------------------------------------------------- |
| LINE_CHANNEL_ACCESS_TOKEN             | LINE Channel Messaging API's channel access token |
| LINE_CHANNEL_SECRET                   | LINE Channel Channel Secret                       |
| GCP_PUBSUB_PROJECT_ID                 | GCP Project ID                                    |
| GCP_PUBSUB_TOPIC                      | Cloud Pub/Sub topic name                          |
| UPLOAD_GOOGLE_DRIVE_ID                | Google Drive ID                                   |
| UPLOAD_SERVICE_ACCOUNT_CLIENT_EMAIL   | service account's client email                    |
| UPLOAD_SERVICE_ACCOUNT_PRIVATE_KEY    | service account's private key                     |
| UPLOAD_SERVICE_ACCOUNT_PRIVATE_KEY_ID | service account's private key id                  |
