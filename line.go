package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	lineServer *linebot.Client
	LINE_ADMIN string
)

func init() {
	var err error
	LINE_SECRET := os.Getenv("LINE_SECRET")
	LINE_TOKEN := os.Getenv("LINE_TOKEN")
	LINE_ADMIN = os.Getenv("LINE_ADMIN")

	lineServer, err = linebot.New(
		LINE_SECRET,
		LINE_TOKEN,
	)
	if err != nil {
		log.Print("Line Service Setup Error: ", err)
		lineServer = nil
	}
}

func LineCallback(w http.ResponseWriter, req *http.Request) {
	log.Println("Callback Request")
	events, err := lineServer.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		profile, err := lineServer.GetProfile(event.Source.UserID).Do()
		if err != nil {
			log.Print("GET PROFILE ERROR: ", err)
			log.Print("PROFILE: ", profile)
			profile = &linebot.UserProfileResponse{}
		}

		log.Print("EventType Received: ", event.Type)
		if event.Type == linebot.EventTypeMessage {

			switch message := event.Message.(type) {
			case *linebot.StickerMessage:
				log.Println("Received Sticker: ", message.StickerID)
				/*
					if _, err = lineServer.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("I have that sticker too!")).Do(); err != nil {
						log.Println("Reply Error: ", err)
					}
					if _, err = lineServer.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.StickerID)).Do(); err != nil {
						log.Println("Reply Error: ", err)
					}*/

			case *linebot.TextMessage:
				log.Printf("Received '%s' from id %s and display name %s", message.Text, profile.UserID, profile.DisplayName)
				var response string

				if strings.ToLower(profile.DisplayName) == "tfan" || strings.ToLower(profile.DisplayName) == "denis" || strings.EqualFold(profile.DisplayName, "TFan") || strings.EqualFold(profile.DisplayName, "Denis") {
					m := strings.ToLower(message.Text)
					if strings.Contains(m, "denis") || strings.Contains(m, "him") || strings.Contains(m, "love") || strings.Contains(m, "too") {
						response = "That's sweet, I will tell my creator for you.  Right. Now."
					} else if strings.Contains(m, "armor") {
						response = "Armor... I heard he is the most handsome dog ever!  Prince!"
					} else {
						response = "Hello Tiffany, my creator loves you more than anything in the world.  Thats what he told me."
					}
				} else {
					response = "text message received"
				}
				log.Println([]byte(profile.DisplayName))
				log.Println([]byte("denis"))
				log.Println(response)

				/*
					if _, err = lineServer.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
						log.Println("Reply Error: ", err)
					}*/

				pub := fmt.Sprintf("Type: \"%s\"  Content: \"%s\"  Response: \"%s\"  FromID: \"%s\"  DisplayName: \"%s\"",
					"text",
					message.Text,
					"none",
					profile.UserID,
					profile.DisplayName)
				subManager.Publish("lnstream", pub)

			case *linebot.ImageMessage:
				/*
					response := "image message received"
					if _, err = lineServer.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
						log.Println("Reply Error: ", err)
					}*/

				log.Println("IMAGE URL: ", message.OriginalContentURL, message.PreviewImageURL)
				pub := fmt.Sprintf("Type: \"%s\"  Content: \"%s\"  Response: \"%s\"  FromID: \"%s\"  DisplayName: \"%s\"",
					"image",
					message.OriginalContentURL,
					"none",
					profile.UserID,
					profile.DisplayName)
				subManager.Publish("lnstream", pub)

			}

		}

		if event.Type == linebot.EventTypeJoin {
			log.Printf("Userid %s - %s joined", profile.UserID, profile.DisplayName)
			/*
				if _, err = lineServer.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Thanks for adding this humble robot, "+profile.DisplayName)).Do(); err != nil {
					log.Println(err)
				}*/
		}
	}
}
