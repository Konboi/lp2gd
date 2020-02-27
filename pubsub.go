package lp2gd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type Message struct {
	ContentIDs []string `json:"content_ids"`
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func HandlePubSub(ctx context.Context, m PubSubMessage) error {
	msg := &Message{}
	if err := json.Unmarshal(m.Data, msg); err != nil {
		log.Println("[ERROR] error unmarshal pubsub message to message", err)
		return err
	}

	email := os.Getenv("CLIENT_EMAIL")
	privateKey := strings.Replace(os.Getenv("PRIVATE_KEY"), `\n`, "\n", -1)
	privateKeyID := os.Getenv("PRIVATE_KEY_ID")
	if email == "" || privateKey == "" || privateKeyID == "" {
		log.Printf("[ERROR] email,private-key, private-key-id are required.")
		return fmt.Errorf("email or private key id or private key is empty")
	}

	scopes := []string{
		drive.DriveScope,
	}

	cfg := &jwt.Config{
		Email:        email,
		PrivateKey:   []byte(privateKey),
		PrivateKeyID: privateKeyID,
		Scopes:       scopes,
		TokenURL:     google.JWTTokenURL,
	}

	client := cfg.Client(ctx)
	dcli, err := drive.New(client)
	if err != nil {
		log.Println("[ERROR] new google drive client", err)
		return err
	}

	for _, contentID := range msg.ContentIDs {
		if err := upload(ctx, dcli, contentID); err != nil {
			log.Println("[ERROR] upload image", err)
			return err
		}
	}

	return nil
}

func upload(ctx context.Context, cli *drive.Service, contentID string) error {
	dir := os.TempDir()
	defer os.RemoveAll(dir)

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Println("[ERROR] new bot client", err)
		return err
	}

	res, err := bot.GetMessageContent(contentID).Do()
	if err != nil {
		log.Println("[ERROR] new request", err)
		return err
	}

	log.Printf("content-id:%s type:%s length:%d\n", contentID, res.ContentType, res.ContentLength)
	f, err := cli.Files.Create(&drive.File{
		Name: fmt.Sprintf("%s.png", contentID),
		Parents: []string{
			os.Getenv("UPLOAD_GOOGLE_DRIVE_ID"),
		},
	}).Media(res.Content, googleapi.ContentType(res.ContentType)).Context(ctx).Do()
	if err != nil {
		log.Println("[ERROR] error post media:", err)
		return err
	}
	fmt.Println("upload success.", "id:", f.Id, "name:", f.Name)

	return nil
}
