package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func exitIfEmpty(value string, message string) {
	if value == "" {
		fmt.Println(message)
		os.Exit(1)
	}
}

func loadConfig() *viper.Viper {
	config := viper.New()
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config.SetDefault("telegram.api", "")
	config.SetDefault("common.name", "")

	config.SetDefault("keel.address", "http://keel.default.svc.cluster.local:9300")
	config.SetDefault("http.listen", "localhost:8000")
	config.SetDefault("telegram.receiver", "498146361")

	exitIfEmpty(config.GetString("telegram.api"), "expects TELEGRAM_API env")
	exitIfEmpty(config.GetString("common.name"), "expects COMMON_NAME env")

	return config
}
