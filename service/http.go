package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type httpInstance struct {
	config *viper.Viper
}

// makes kubernetes happy
func (h *httpInstance) health(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK!\n"))
}

func initHTTPRouter(config *viper.Viper) *mux.Router {
	router := mux.NewRouter()
	http := &httpInstance{}

	http.config = config
	router.HandleFunc("/health", http.health).Methods("GET", "OPTIONS")

	return router
}
