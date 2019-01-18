package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (s *Sevice) updateKeelDeployment(image string) error {
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return errors.New("expect {name:tag} as image")
	}
	name, tag := parts[0], parts[1]

	// r := bytes.NewReader(data)
	url := fmt.Sprintf("%s/v1/webhooks/native", s.config.GetString("keel.address"))

	values := map[string]string{"name": name, "tag": tag}
	jsonValue, _ := json.Marshal(values)
	log.Println(values)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	return err
}
