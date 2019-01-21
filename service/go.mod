module mod

replace github.com/afoninsky/makeomatic/common => ./common

replace github.com/afoninsky/makeomatic/providers/telegram => ./providers/telegram

require (
	github.com/afoninsky/makeomatic/common v0.0.0
	github.com/afoninsky/makeomatic/providers/telegram v0.0.0
	github.com/gorilla/mux v1.6.2
	github.com/spf13/viper v1.3.1
)
