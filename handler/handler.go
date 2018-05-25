package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	lineGreetingMessage = `Hi! Kentang's here. Add this bot to your group and count your friends koplaquness!`
	lineHelpString      = `Here are available commands:
add [keyword] [desc]
remove [keyword]
list
[keyword] -> Increase count
highschore -> This month only
statistics ->
reset -> Reset all
help -> Show this`
)

type Handler struct {
	bot *linebot.Client
}

func NewHandler(channelSecret, channelToken string) *Handler {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &Handler{bot: bot}
}

//Healthz - health check
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	log.Printf("health\n")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}

func (h *Handler) log(format string, args ...interface{}) {
	log.Printf("[LINE] "+format, args...)
}

func (h *Handler) reply(event *linebot.Event, messages ...string) error {
	var lineMessages []linebot.Message
	for _, message := range messages {
		lineMessages = append(lineMessages, linebot.NewTextMessage(message))
	}
	_, err := h.bot.ReplyMessage(event.ReplyToken, lineMessages...).Do()
	if err != nil {
		h.log("Error replying to %+v: %s", event.Source, err.Error())
	}
	return err
}

func (h *Handler) push(to string, messages ...string) error {
	var lineMessages []linebot.Message
	for _, message := range messages {
		lineMessages = append(lineMessages, linebot.NewTextMessage(message))
	}
	_, err := h.bot.PushMessage(to, lineMessages...).Do()
	if err != nil {
		h.log("Error pushing to %s: %s", to, err.Error())
	}
	return err
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := h.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		h.log("[EVENT][%s] Source: %#v", event.Type, event.Source)
		switch event.Type {

		case linebot.EventTypeJoin:
			fallthrough
		case linebot.EventTypeFollow:
			h.handleFollow(event)

		case linebot.EventTypeLeave:
			fallthrough
		case linebot.EventTypeUnfollow:
			fallthrough

		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				h.handleTextMessage(event, message)
			}
		}
	}
}

func (h *Handler) handleFollow(event *linebot.Event) {
	message := lineGreetingMessage + "\n\n" + lineHelpString
	h.reply(event, message)
}

func (h *Handler) handleTextMessage(event *linebot.Event, message *linebot.TextMessage) {

}
