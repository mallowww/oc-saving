package main

import (
	"bytes"
	"encoding/json"
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

	if configFile, err := os.ReadFile("./config/config.yml"); err != nil {
		log.Fatalf("error from reading config file: %v", err)
	} else {
		expandedConfig := os.ExpandEnv(string(configFile))
		viper.ReadConfig(bytes.NewBufferString(expandedConfig))
	}

	log.Println("server start at port:", viper.GetInt("app.port"))

	// log
	if configJSON, err := json.MarshalIndent(viper.AllSettings(), "", "  "); err != nil {
		log.Fatalf("error from marshalling config: %v", err)
	} else {
		log.Printf("config:\n%s", string(configJSON))
	}
}
