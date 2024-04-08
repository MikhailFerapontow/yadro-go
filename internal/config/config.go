package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func MustLoad(config_path string) {
	fmt.Println("config_path", config_path)
	viper.AddConfigPath(config_path)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Config initialization failed with %s", err)
	}
}
