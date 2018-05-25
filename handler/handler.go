package handler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/luqmanarifin/kentang/model"
	"github.com/luqmanarifin/kentang/service"
	"github.com/luqmanarifin/kentang/util"
)

var (
	lineGreetingMessage = `Hi! Kentang's here. Add this bot to your group and count your friends koplaqueness!`
	lineHelpString      = `Here are available commands:
- add [keyword] [desc]
- remove [keyword]
- list
- [keyword] -> Increase count
- highscore -> This month
- stat
- reset -> Reset all
- help`
)

type Handler struct {
	bot   *linebot.Client
	mysql *service.MySQL
}

func NewHandler(channelSecret, channelToken string) *Handler {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	opt := service.MySQLOption{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Database: os.Getenv("MYSQL_DATABASE"),
		Charset:  os.Getenv("MYSQL_CHARSET"),
	}
	mysql, err := service.NewMySQL(opt)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return &Handler{bot: bot, mysql: mysql}
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
	source := util.LineEventSourceToReplyString(event.Source)
	log.Printf("Received message from %s: %s", source, message.Text)

	tokens := strings.Split(message.Text, " ")
	switch tokens[0] {
	case "add":
		h.handleAdd(event, tokens)
	case "remove":
		h.handleRemove(event, tokens)
	case "list":
		h.handleList(event, tokens)
	case "highscore":
		h.handleHighscore(event, tokens)
	case "stat":
		h.handleStat(event, tokens)
	case "reset":
		h.handleReset(event, tokens)
	case "help":
		h.handleHelp(event, tokens)
	case "profile":
		h.handleProfile(event, tokens)
	default:
		h.handleKeyword(event, tokens)
	}
}

func (h *Handler) handleAdd(event *linebot.Event, tokens []string) {
	if len(tokens) != 3 {
		return
	}
	keyword := tokens[1]
	desc := tokens[2]
	source := util.LineEventSourceToReplyString(event.Source)

	dict, err := h.mysql.GetDictionaryByKeyword(source, keyword)
	if err != nil {
		log.Printf("Error when getting info when adding %s in %s\n", keyword, source)
		return
	}
	if dict.Keyword == keyword {
		h.reply(event, "Keyword "+keyword+" already registered here before.")
		return
	}
	err = h.mysql.CreateDictionary(&model.Dictionary{
		Source:      source,
		Keyword:     keyword,
		Description: desc,
	})
	if err != nil {
		log.Printf("Error when adding %s in %s\n", keyword, source)
		return
	}
	h.reply(event, "Keyword "+keyword+" has been added")
}

func (h *Handler) handleRemove(event *linebot.Event, tokens []string) {
	if len(tokens) != 2 {
		return
	}
	keyword := tokens[1]
	source := util.LineEventSourceToReplyString(event.Source)

	dict, err := h.mysql.GetDictionaryByKeyword(source, keyword)
	if err != nil {
		log.Printf("Error when getting info when deleting %s in %s\n", keyword, source)
		return
	}
	if dict.Keyword != keyword {
		h.reply(event, "Keyword "+keyword+" hasn't been registered")
		return
	}
	err = h.mysql.RemoveDictionary(&dict)
	if err != nil {
		log.Printf("Error when deleting %s in %s\n", keyword, source)
		return
	}
	h.reply(event, "Keyword "+keyword+" hasn't been removed")
}

func (h *Handler) handleList(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
	source := util.LineEventSourceToReplyString(event.Source)
	dicts, err := h.mysql.GetAllDictionaries(source)
	if err != nil {
		log.Printf("Error when fetching dictionaries for %s\n", source)
		return
	}
	if len(dicts) == 0 {
		h.reply(event, "No keyword registered at this time.")
		return
	}
	message := "Keywords registered in this group:"
	for _, dict := range dicts {
		message = message + "\n- " + dict.Keyword + ": " + dict.Description
	}
	h.reply(event, message)
}

func (h *Handler) handleHighscore(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
	source := util.LineEventSourceToReplyString(event.Source)
	entries, err := h.mysql.GetMonthEntries(source)
	if err != nil {
		log.Printf("Error in fetching highscore")
		return
	}
	m := util.EntriesToSortedMap(entries)
	message := "Highscore for this month:"
	for key, value := range m {
		message = message + "\n" + key + ": " + strconv.Itoa(value)
	}
	h.reply(event, message)
}

func (h *Handler) handleStat(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
}

func (h *Handler) handleReset(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
	source := util.LineEventSourceToReplyString(event.Source)
	err := h.mysql.RemoveDictionaryBySource(source)
	if err != nil {
		log.Printf("Error in resetting dictionary in %s", source)
		return
	}
	err = h.mysql.RemoveEntryBySource(source)
	if err != nil {
		log.Printf("Error in resetting source in %s", source)
		return
	}
	h.reply(event, "Has been reset")
}

func (h *Handler) handleHelp(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
	h.reply(event, lineHelpString)
}

func (h *Handler) handleProfile(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
	source := event.Source
	if source.UserID != "" {
		profile, err := h.bot.GetProfile(source.UserID).Do()
		if err != nil {
			h.reply(event, err.Error())
			return
		}
		if _, err := h.bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage("Display name: "+profile.DisplayName),
			linebot.NewTextMessage("Status message: "+profile.StatusMessage),
		).Do(); err != nil {
			return
		}
	} else {
		h.reply(event, "Bot can't use profile API without user ID")
	}
}

func (h *Handler) handleKeyword(event *linebot.Event, tokens []string) {
	if len(tokens) != 1 {
		return
	}
	keyword := tokens[0]
	source := util.LineEventSourceToReplyString(event.Source)

	// TODO: find keyword on redis first

	// find keyword on mysql
	dict, err := h.mysql.GetDictionaryByKeyword(source, keyword)
	if err != nil || dict.Keyword != keyword {
		log.Printf("Can't found keyword %s in %s\n", keyword, source)
		return
	}
	err = h.mysql.CreateEntry(&model.Entry{
		Keyword: keyword,
		Source:  source,
	})
	if err != nil {
		log.Printf("Cannot add counter %s in %s\n", keyword, source)
		return
	}
	h.reply(event, "Keyword "+keyword+" added.")
}
