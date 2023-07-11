package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Load the config.yml file and watch it for any updates
func loadConfig() error {
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		newErr := fmt.Sprintf(
			"CONFIG ERR: %v\n\nCopy `config.example.yml` to `config.yml`\nAnd adjust `config.yml` to your liking", err)
		return errors.New(newErr)
	}
	viper.WatchConfig()
	return nil
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

func getMaxLat() int {
	return viper.GetInt("server.max_latency")
}
func getPingerTimeout() int {
	return viper.GetInt("server.pinger_timeout")
}
