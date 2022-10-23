package config

import "os"

type Config struct {
	BitFlyerApiKey    string
	BitFlyerApiSecret string
}

func NewConfig() Config {
	return Config{
		BitFlyerApiKey:    os.Getenv("BITFLYER_API_KEY"),
		BitFlyerApiSecret: os.Getenv("BITFLYER_API_SECRET"),
	}
}
