package config

import "github.com/spf13/viper"

// Config contains all needed data for notification processing
type Config struct {
	KafkaTopic    string   `mapstructure:"KAFKA_TOPIC"`
	SmtpFrom      string   `mapstructure:"SMTP_FROM"`
	ConsumerGroup string   `mapstructure:"CONSUMER_GROUP"`
	SmtpHost      string   `mapstructure:"SMTP_HOST"`
	SmtpPort      string   `mapstructure:"SMTP_PORT"`
	SmtpUsername  string   `mapstructure:"SMTP_USERNAME"`
	SmtpPassword  string   `mapstructure:"SMTP_PASSWORD"`
	KafkaBrokers  []string `mapstructure:"KAFKA_BROKERS"`
}

// LoadConfig loads the app.env file to the Config struct.
// If there are any ENV variables with the same name as Config mapstruct field tags
// this data is overwritten by environment variables
func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
