package main

import (
	"log"
	"net/http"
	"os"

	"github.com/s14t284/every-news-bot/utils"
)

func main() {
	handler, err := utils.GetLineBotHandler()
	if err != nil {
		return
	}
	http.Handle("/callback", handler)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
