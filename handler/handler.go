package handler

import (
	"net/http"
)

type Handler struct {
	// bot *linebot.Client
}

func NewHandler(channelSecret, channelToken string) *Handler {
	// bot, err := linebot.New(
	// 	os.Getenv("CHANNEL_SECRET"),
	// 	os.Getenv("CHANNEL_TOKEN"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return &Handler{}
}

//Healthz - health check
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}

//Index
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}
