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
	config.SetConfigName("makeobot")

	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	return config
}
