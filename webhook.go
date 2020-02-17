package lp2gp

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

const ContentEndpoint = "https://api-data.line.me/v2/bot/message/%s/content"

type Request struct {
	Events []*linebot.Event `json:"events"`
}

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	secret := os.Getenv("CHANNEL_SECRET")
	if secret == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("body:%s", string(body))

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	events, err := linebot.ParseRequest(secret, r)
	if err != nil {
		log.Printf("[ERROR] web hook json encode faild. error:%v\nsecret:%s, header:%s", err, secret, r.Header.Get("X-Line-Signature"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, e := range events {
		if e.Type != linebot.EventTypeMessage {
			continue
		}

		switch msg := e.Message.(type) {
		case *linebot.TextMessage:
			log.Printf("[INFO] text message %s", msg.Text)
		case *linebot.ImageMessage:
			log.Println("--- content ---")
			log.Printf(ContentEndpoint, msg.ID)
			log.Printf("[INFO] image message. id:%s original:%s preview:%s\nmsg:%+v", msg.ID, msg.OriginalContentURL, msg.PreviewImageURL, msg)
		default:
			log.Printf("[ERROR] unsupport message type. type:%T", e.Message)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		log.Printf("[INFO] web hook message. eventt: %+v", e)
	}
	return
}
