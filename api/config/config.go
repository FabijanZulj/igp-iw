package config

import "github.com/spf13/viper"

// Config contains all needed data for the api
type Config struct {
	KafkaTopic   string   `mapstructure:"KAFKA_TOPIC"`
	DBSource     string   `mapstructure:"DB_SOURCE"`
	JwtSecret    string   `mapstructure:"JWT_SECRET"`
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
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
