package handler

import (
	"log"
	"net/http"
	"net/url"
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
	redis *service.Redis
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

	redisUrl, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	password, _ := redisUrl.User.Password()
	redisOpt := service.RedisOption{
		Host:     redisUrl.Hostname(),
		Port:     redisUrl.Port(),
		Password: password,
		Database: 0,
	}
	log.Printf("redis opt: %v", redisOpt)
	redis, err := service.NewRedis(redisOpt)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return &Handler{bot: bot, mysql: mysql, redis: redis}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cok"))
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
	switch strings.ToLower(tokens[0]) {
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

	val, err := h.redis.GetKeyword(source, keyword)

	if err == nil && val != util.NOT_EXIST {
		h.reply(event, keyword+" is already here before.")
		return
	}

	dict, err := h.mysql.GetDictionaryByKeyword(source, keyword)
	if dict.Keyword == keyword {
		h.reply(event, keyword+" is already here before.")
		return
	}
	err = h.mysql.CreateDictionary(&model.Dictionary{
		Source:      source,
		Keyword:     keyword,
		Description: desc,
		Creator:     event.Source.UserID,
	})
	if err != nil {
		log.Printf("Error when adding %s in %s\n", keyword, source)
		return
	}
	h.reply(event, keyword+" has been added")

	err = h.redis.AddKeyword(source, keyword, desc)
	if err != nil {
		log.Printf("Error when adding cache %s in %s\n", keyword, source)
		return
	}
}

func (h *Handler) handleRemove(event *linebot.Event, tokens []string) {
	if len(tokens) != 2 {
		return
	}
	keyword := tokens[1]
	source := util.LineEventSourceToReplyString(event.Source)

	desc, err := h.redis.GetKeyword(source, keyword)
	if err == nil && desc == util.NOT_EXIST {
		h.reply(event, "Keyword "+keyword+" is not exists")
		return
	}

	dict, err := h.mysql.GetDictionaryByKeyword(source, keyword)
	if err != nil {
		log.Printf("Error when getting info when deleting %s in %s\n", keyword, source)
		return
	}
	if dict.Keyword != keyword {
		h.reply(event, "Keyword "+keyword+" is not exists")
		return
	}
	if event.Source.UserID != dict.Creator {
		h.reply(event, "Only the creator can remove it")
		return
	}
	err = h.mysql.RemoveDictionary(&dict)
	if err != nil {
		log.Printf("Error when deleting %s in %s\n", keyword, source)
		return
	}
	err = h.mysql.RemoveEntryByKeyword(source, keyword)
	if err != nil {
		log.Printf("Error when deleting entries %s in %s", keyword, source)
		return
	}
	h.reply(event, "Keyword "+keyword+" removed")

	err = h.redis.RemoveKeyword(source, keyword)
	if err != nil {
		log.Printf("Error when deleting cache %s in %s\n", keyword, source)
	}
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
		h.reply(event, "No keyword registered.")
		return
	}
	message := "Keywords:"
	for i, dict := range dicts {
		message = message + "\n" + strconv.Itoa(i+1) + ". " + dict.Keyword + ": " + dict.Description + " (" + h.getProfileName(dict.Creator) + ")"
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
	p := util.EntriesToSortedMap(entries)
	if len(p) == 0 {
		h.reply(event, "No highscore")
		return
	}
	message := "Highscore:"
	for _, pair := range p {
		dict, err := h.mysql.GetDictionaryByKeyword(source, pair.Value)
		if err != nil {
			log.Printf("error when fetching desc")
			return
		}
		message = message + "\n" + pair.Value + " - " + dict.Description + " : " + strconv.Itoa(pair.Key)
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
	h.reply(event, "All cleared up.")

	err = h.redis.RemoveAllKeyword(source)
	if err != nil {
		log.Printf("Error in resetting cache source in %s", source)
		return
	}
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

	// find keyword on redis first
	ret, err := h.redis.GetKeyword(source, keyword)
	if ret == util.NOT_EXIST {
		log.Printf("%s is NOT exist in %s, based on cache", keyword, source)
		return
	} else if err == nil {
		log.Printf("%s is exist in %s, based on cache", keyword, source)
		h.addEntry(event, source, keyword, ret)
		return
	}

	// find keyword on mysql
	dict, err := h.mysql.GetDictionaryByKeyword(source, keyword)
	if err != nil || dict.Keyword != keyword {
		h.redis.RemoveKeyword(source, keyword)
		log.Printf("Can't found keyword %s in %s, found %s\n", keyword, source, dict.Keyword)
		return
	}
	h.redis.AddKeyword(source, keyword, dict.Description)
	h.addEntry(event, source, keyword, dict.Description)
}

func (h *Handler) addEntry(event *linebot.Event, source, keyword, desc string) {
	err := h.mysql.CreateEntry(&model.Entry{
		Keyword: keyword,
		Source:  source,
	})
	if err != nil {
		log.Printf("Cannot add counter %s in %s\n", keyword, source)
		return
	}
	h.reply(event, keyword+", "+desc+" lagi?")
}

func (h *Handler) getProfileName(userId string) string {
	// look up from cache
	name, _ := h.redis.GetDisplayName(userId)
	if name != "" {
		return name
	}

	// look up from database
	profile, err := h.bot.GetProfile(userId).Do()
	if err != nil {
		return ""
	}
	// update cache
	h.redis.SetDisplayName(userId, profile.DisplayName)
	return profile.DisplayName
}
