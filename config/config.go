package config

import "os"

type Config struct {
	BitFlyerApiKey    string
	BitFlyerApiSecret string
	KabucomAPIHost    string
}

func NewConfig() Config {
	return Config{
		/* BitFlyer */
		BitFlyerApiKey:    os.Getenv("BITFLYER_API_KEY"),
		BitFlyerApiSecret: os.Getenv("BITFLYER_API_SECRET"),
		/* Kabucom */
		KabucomAPIHost: os.Getenv("KABUCOM_API_HOST"),
	}
}
