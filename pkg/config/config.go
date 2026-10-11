package config

import "github.com/spf13/viper"

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	CORS    CORSConfig    `mapstructure:"cors"`
	Anchors []Anchor      `mapstructure:"anchors"`
	Goal    Goal          `mapstructure:"goal"`
}

type ServerConfig struct {
	HTTP HTTPConfig `mapstructure:"http"`
}

type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type Anchor struct {
	ID string  `mapstructure:"id"`
	X  float64 `mapstructure:"x"`
	Y  float64 `mapstructure:"y"`
	Z  float64 `mapstructure:"z"`
}

type Goal struct {
	X    float64 `mapstructure:"x"`
	YMin float64 `mapstructure:"y_min"`
	YMax float64 `mapstructure:"y_max"`
}

func Load() (Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
