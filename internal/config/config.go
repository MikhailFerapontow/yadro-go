package config

import (
	"log"

	"github.com/spf13/viper"
)

func MustLoad() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Config initialization failed with %s", err)
	}
}
