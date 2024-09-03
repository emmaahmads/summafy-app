package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Port          string `mapstructure:"PORT"`
	DBUrl         string `mapstructure:"DB_URL"`
	DBDriver      string `mapstructure:"DB_DRIVER"`
	S3Bucket      string `mapstructure:"S3_BUCKET"`
	Region        string `mapstructure:"AWS_REGION"`
	Creds1        string `mapstructure:"AWS_CREDENTIALS_1"`
	Creds2        string `mapstructure:"AWS_CREDENTIALS_2"`
	Creds3        string `mapstructure:"AWS_CREDENTIALS_3"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
