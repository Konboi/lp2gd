package lp2gp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/line/line-bot-sdk-go/linebot"
)

const ContentEndpoint = "https://api-data.line.me/v2/bot/message/%s/content"

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	secret := os.Getenv("LINE_CHANNEL_SECRET")
	if secret == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projectID := os.Getenv("GCP_PUBSUB_PROJECT_ID")
	if projectID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	topic := os.Getenv("GCP_PUBSUB_TOPIC")
	if topic == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := linebot.ParseRequest(secret, r)
	if err != nil {
		log.Printf("[ERROR] web hook json encode faild. error:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := []string{}
	for _, e := range events {
		if e.Type != linebot.EventTypeMessage {
			continue
		}

		switch msg := e.Message.(type) {
		case *linebot.ImageMessage:
			ids = append(ids, msg.ID)
		default:
			continue
		}
	}

	if len(ids) == 0 {
		fmt.Fprintf(w, "OK")
		return
	}

	msg := Message{
		ContentIDs: ids,
	}

	client, err := pubsub.NewClient(r.Context(), projectID)
	if err != nil {
		log.Printf("[ERROR] error new pubsub client. err:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[ERROR] error mershal message. err:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := client.Topic(topic)
	t.Publish(r.Context(), &pubsub.Message{Data: data})
}
