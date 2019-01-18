package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type keelDeploymentData struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

// const keelUpdateTemplate = template.Must(template.New("keel-update").Parse("*deployment update*\nSuccessfully updated bla-bla"))

// HealthHandler makes k8s be happy
func (s *Sevice) httpHealthHandler(res http.ResponseWriter, req *http.Request) {
	// TODO: check keel, ex,: curl http://keel.default.svc.cluster.local:9300/healthz
	res.Write([]byte("OK!\n"))
}

// KeelDeploymentHandler receives deployment information from keel and
// sends it to the telegram
func (s *Sevice) httpKeelNewDeploymentHandler(res http.ResponseWriter, req *http.Request) {
	var data keelDeploymentData
	dec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if err := dec.Decode(&data); err != nil {
		http.Error(res, "unable to decode body", http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer
	if err := s.tplNewDeployment.Execute(&buffer, data); err != nil {
		http.Error(res, "unable to parse remplate", http.StatusBadRequest)
		return
	}
	msg := tgbotapi.NewMessageToChannel(s.config.GetString("telegram.receiver"), buffer.String())
	msg.ParseMode = "Markdown"

	s.bot.Send(msg)

	res.Write([]byte("OK!\n"))
}
