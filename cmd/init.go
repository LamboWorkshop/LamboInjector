package cmd

import (
	"io"
	"log"
	"os"
	"time"

	api "github.com/segfault42/binance-api"
	"github.com/spf13/viper"
)

type SConfig struct {
	LogFilePath string `json:"logFilePath"`
	TimeZone    string `json:"timeZone"`
}

func initLoging(logFilePath, timeZone string) (*os.File, error) {
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Set output to file and stdout
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	// get the location
	location, _ := time.LoadLocation(timeZone)
	log.Println("======= New Log Session : " + time.Now().In(location).String() + " =======\n")

	return f, err
}

func loadConfig() (SConfig, error) {

	var config SConfig
	var err error

	if os.Getenv("APP_ENV") == "production" {
		viper.SetConfigFile("config_prod.json")
	} else {
		viper.SetConfigFile("config_dev.json")
	}

	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return SConfig{}, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return SConfig{}, err
	}

	return config, nil
}

func InitServices() (*os.File, api.ApiInfo, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, api.ApiInfo{}, err
	}

	f, err := initLoging(config.LogFilePath, config.TimeZone)
	if err != nil {
		return nil, api.ApiInfo{}, err
	}

	return f, api.New(), nil
}
