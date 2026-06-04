package config

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	MinConns int    `mapstructure:"min_conns"`
	MaxConns int    `mapstructure:"max_conns"`
}

func MustLoadConfig(path string) Config {
	viper.AddConfigPath(path)
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	viper.SetDefault("database.min_conns", 1)
	viper.SetDefault("database.max_conns", 1)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	return config
}
