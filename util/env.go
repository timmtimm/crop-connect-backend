package util

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func GetConfig(key string) string {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error when reading config: %s", err)
	}

	return viper.GetString(key)
}

func ResontructeDomainName() []string {
	return strings.Split(GetConfig("DOMAIN_NAME"), ",")
}
