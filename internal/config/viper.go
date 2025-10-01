package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, continue with system env")
	}

	config := viper.New()

	config.AutomaticEnv()

	return config
}
