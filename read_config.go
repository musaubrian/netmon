package main

import (
	"log"

	"github.com/spf13/viper"
)

// Load the config.yml file and watch it for any updates
func loadConfig() {
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("CONFIG ERR: ", err,
			"\n\nCopy `config.example.yml` to `config.yml`\nAnd adjust `config.yml` to your liking")
	}
	viper.WatchConfig()
}

func getServerToPing() string {
	return viper.GetString("server.server_to_ping")
}

func getPort() int {
	return viper.GetInt("server.web_server_port")
}

func getEmails() []string {
	return viper.GetStringSlice("emails")
}
