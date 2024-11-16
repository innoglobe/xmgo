package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Production bool
	Server     ServerConf
	Database   DBConf
	JWT        JWTConf
	Kafka      KafkaConfig
}

type ServerConf struct {
	Timeout int
	Port    int
	Host    string
	SSL     SSLConf
}

type SSLConf struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type DBConf struct {
	Port int
	Host string
	User string
	Pass string
	Name string
}

type JWTConf struct {
	Secret string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func LoadConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
