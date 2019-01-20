package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	config := loadConfig()
	router := initHTTPRouter(config)

	// TODO: load services

	// TODO: load notificators

	fmt.Printf("start http service on %s\n", config.GetString("http.listen"))
	log.Fatal(http.ListenAndServe(config.GetString("http.listen"), router))

}
