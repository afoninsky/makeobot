package main

import (
	"net/http"
	"github.com/afoninsky/makeomatic/providers/telegram"
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

	if err := router.RegisterService("telegram", ctx, &telegram.Service{}); err != nil {
		logger.Fatal(err)
	}

	logger.Printf("listen on %s\n", config.GetString("http.listen"))
	logger.Fatal(http.ListenAndServe(config.GetString("http.listen"), httpMux))

}
