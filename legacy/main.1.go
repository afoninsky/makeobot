package main

import (
	"fmt"
	"log"
	"net/http"
)

func main1() {

	config := loadConfig()
	bot, tgErr := initTgBot(config)
	if tgErr != nil {
		log.Fatal(tgErr)
	}
	router := initHTTPRouter(config, bot)
	fmt.Printf("start http service on %s\n", config.GetString("http.listen"))
	log.Fatal(http.ListenAndServe(config.GetString("http.listen"), router))
}
