package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/luqmanarifin/kentang/handler"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	log.Printf("channel secret %s\n", os.Getenv("CHANNEL_SECRET"))
	log.Printf("channel token %s\n", os.Getenv("CHANNEL_TOKEN"))
	log.Printf("port %s\n", os.Getenv("PORT"))

	handler := handler.NewHandler(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/healthz", handler.Healthz)
	http.HandleFunc("/callback", handler.Callback)

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
