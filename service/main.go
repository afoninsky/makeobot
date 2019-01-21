package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	_ "github.com/afoninsky/makeomatic/providers/telegram"
)

type instance struct {
	config *viper.Viper
	router *mux.Router
}

func main() {

	config := loadConfig()
	router := initHTTPRouter(config)

	// service := &instance{
	// 	config: config,
	// 	router: router,
	// }
	// TODO: load services

	// TODO: load notificators

	fmt.Printf("start http service on %s\n", config.GetString("http.listen"))
	log.Fatal(http.ListenAndServe(config.GetString("http.listen"), router))

}
