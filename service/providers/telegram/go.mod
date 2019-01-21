module github.com/afoninsky/makeomatic/providers/telegram

replace github.com/afoninsky/makeomatic/common => ../../common

require (
	github.com/afoninsky/makeomatic/common v0.0.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/gorilla/mux v1.6.2
	github.com/spf13/viper v1.3.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
