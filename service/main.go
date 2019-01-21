package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/afoninsky/makeomatic/providers/telegram"
	"github.com/afoninsky/makeomatic/common"
)

type instance struct {
	config *viper.Viper
	router *mux.Router
}

func main() {

	var notifier common.Notifier

	config := loadConfig()
	router := initHTTPRouter(config)

	// service := &instance{
	// 	config: config,
	// 	router: router,
	// }

	// init notification service
	notifier, nErr := telegram.New(config, router)
	if nErr != nil {
		log.Fatal(nErr)
	}
	defer notifier.Close()


	fmt.Printf("start http service on %s\n", config.GetString("http.listen"))
	log.Fatal(http.ListenAndServe(config.GetString("http.listen"), router))

}
