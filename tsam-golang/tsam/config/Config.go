package config

import (
	"github.com/spf13/viper"
	log "github.com/techlabs/swabhav/tsam/log"
)

// Config Contain Viper
type Config struct {
	viper ConfReader
}

// ConfReader defines all methods to be present in Config.
type ConfReader interface {
	GetString(key string) string
	IsSet(key string) bool
	GetInt64(key string) int64
}

// NewConfig Read envfile and Return Config
func NewConfig(isProduction bool) ConfReader {

	log.NewLogger().Print(" ================= isProduction ->", isProduction)
	vp := viper.New()
	if isProduction {
		vp.SetConfigName("config")
	} else {
		vp.SetConfigName("config-local")
	}
	vp.SetConfigType("env")
	vp.AddConfigPath(".")
	vp.AutomaticEnv()

	config := Config{
		viper: vp,
	}

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.NewLogger().Warn("file Not Found")
		} else {
			log.NewLogger().Fatalf("Something Wrong in File Reading Error:[%s]", err.Error())
		}
	}
	return &config
}

// GetString Return Environment Variable in string type
func (config *Config) GetString(key string) string {
	return config.viper.GetString(key)
}

// IsSet Check Environment Variable Set Or Not.
func (config *Config) IsSet(key string) bool {
	return config.viper.IsSet(key)
}

// GetInt64 Return Environment Variable in int64 type
func (config *Config) GetInt64(key string) int64 {
	return config.viper.GetInt64(key)
}
