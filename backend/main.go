package main

import (
	"bytes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	// env
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("error from reading .env file")
	}
	// config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// viper.AutomaticEnv()
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("error from reading config file: %v", err)
	// }

	// configData := viper.AllSettings()
	configFile, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalf("error from reading config file: %v", err)
	}
	expandedConfig := os.ExpandEnv(string(configFile))

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBufferString(expandedConfig))

	// var cfg config.AppConfig
	// if err := viper.Unmarshal(&cfg); err != nil {
	// 	log.Fatalf("error from unmarshalling config: %v", err)
	// }
	// log.Printf("config loaded: %+v", cfg)
	log.Println("server start")
	log.Printf("app name: %s", viper.GetString("app.name"))
	log.Printf("config: %+v", viper.AllSettings())
}
