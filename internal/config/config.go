package config

import (
	"encoding/json"
	"log"
	"os"
)

type BotConfig struct {
	Token string `json:"token"`
}

// type PostgresConfig struct {
// 	Host    string `json:"host"`
// 	Port    string `json:"port"`
// 	User    string `json:"user"`
// 	DBName  string `json:"dbname"`
// 	SSLMode string `json:"sslmode"`
// }

type Config struct {
	Bot      BotConfig      `json:"bot"`
	// Postgres PostgresConfig `json:"postgres"`
}

func LoadConfig() Config {
	var config Config
	data, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		log.Fatalf("cannot load bot config %v", err)
	}
	if err = json.Unmarshal(data, &config); err != nil {
		log.Fatalf("cannot load bot config %v", err)
	}

	return config
}