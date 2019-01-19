package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	config := loadConfig()
	bot, tgErr := initTgBot(config)
	if tgErr != nil {
		log.Fatal(tgErr)
	}
	router := initHTTPRouter(config, bot)
	fmt.Printf("[%s] start http service on %s\n", config.GetString("common.name"), config.GetString("http.listen"))
	log.Fatal(http.ListenAndServe(config.GetString("http.listen"), router))
}
