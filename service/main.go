package main

import (
	"net/http"
	"github.com/afoninsky/makeomatic/providers/telegram"
	"github.com/afoninsky/makeomatic/providers/keel"
	"github.com/afoninsky/makeomatic/common"
)


func main() {

	logger := common.CreateLogger("core")
	config := loadConfig()
	httpMux := initHTTP(config)
	router := InitServiceRouter()

	ctx := &common.AppContext{
		Config: config,
		HTTP: httpMux,
		Router: router,
	}

	services := map[string]common.ServiceProvider {
		"telegram": &telegram.Service{},
		"keel": &keel.Service{},
	}

	for name, service := range services {
		if err := router.RegisterService(name, ctx, service); err != nil {
			logger.Fatal(err)
		}
	}
	

	logger.Printf("listen on %s\n", config.GetString("http.listen"))
	logger.Fatal(http.ListenAndServe(config.GetString("http.listen"), httpMux))
}
