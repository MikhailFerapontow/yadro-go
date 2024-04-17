package config

import (
	"log"

	"github.com/spf13/viper"
)

func MustLoad(config_path string) {
	viper.AddConfigPath(config_path)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Config initialization failed with %s", err)
	}
	log.Println("Config loaded")
}
