package main

import (
	"log"

	"github.com/spf13/viper"
)

// Load the config.yaml file and update the config values change
func loadConfig() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("CONFIG ERR: ", err,
			"\n\nCopy `config.example.yaml` to `config.yaml`\nAnd adjust `config.yaml` to your liking")
	}
	viper.WatchConfig()
}

func getServerToPing() string {
	loadConfig()
	return viper.GetString("server.server_to_ping")
}

func getPort() int {
	loadConfig()
	return viper.GetInt("server.web_server_port")
}

func getEmails() []string {
	loadConfig()
	return viper.GetStringSlice("emails")
}
