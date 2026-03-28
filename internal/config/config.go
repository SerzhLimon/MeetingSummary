package config

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

type BotConfig struct {
	Token string `json:"token"`
}

type SaluteSpeech struct {
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
}

type GigaChat struct {
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Password string `json:"password"`
}

type Config struct {
	Bot            BotConfig      `json:"bot"`
	Postgres       PostgresConfig `json:"postgres"`
	Salute         SaluteSpeech   `json:"salute_speech"`
	GigaChat       GigaChat       `json:"giga_chat"`
	IntervalTicker int            `json:"interval_ticker"`
}

func LoadConfig() *Config {
	var config Config
	data, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		logrus.Fatalf("cannot load bot config %v", err)
	}
	if err = json.Unmarshal(data, &config); err != nil {
		logrus.Fatalf("cannot load bot config %v", err)
	}
	logrus.Info("config create success")
	return &config
}
