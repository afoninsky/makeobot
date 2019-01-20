// hooks for external services (keel.sh, etc)
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	tpl "text/template"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type keelDeploymentData struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Service string
}

type httpInstance struct {
	config             *viper.Viper
	bot                *tgbotapi.BotAPI
	tplDeploymentEvent *tpl.Template
}

// makes kubernetes happy
func (h *httpInstance) health(res http.ResponseWriter, req *http.Request) {
	// TODO: check keel, ex,: curl http://keel.default.svc.cluster.local:9300/healthz
	res.Write([]byte("OK!\n"))
}

// handle triggering of the new deployment by keel.sh
func (h *httpInstance) keelDeploymentEvent(res http.ResponseWriter, req *http.Request) {
	// decode incoming request
	var data keelDeploymentData
	dec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if err := dec.Decode(&data); err != nil {
		http.Error(res, "unable to decode body", http.StatusBadRequest)
		return
	}

	// send data to the messenger
	var buffer bytes.Buffer
	if err := h.tplDeploymentEvent.Execute(&buffer, data); err != nil {
		http.Error(res, "unable to parse remplate", http.StatusBadRequest)
		return
	}
	msg := tgbotapi.NewMessageToChannel(h.config.GetString("telegram.receiver"), buffer.String())
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
	res.Write([]byte("OK!\n"))
}

func initHTTPRouter(config *viper.Viper, bot *tgbotapi.BotAPI) *mux.Router {
	router := mux.NewRouter()
	http := &httpInstance{}

	http.config = config
	http.bot = bot
	http.tplDeploymentEvent = tpl.Must(tpl.New("keel-deployment").Parse(deploymentEventTemplate))

	router.HandleFunc("/health", http.health).Methods("GET", "OPTIONS")

	router.HandleFunc("/hook/keel/deployment", http.keelDeploymentEvent).Methods("POST", "OPTIONS")

	return router
}
