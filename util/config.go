package util

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	Port           string `mapstructure:"PORT"`
	DBUrl          string `mapstructure:"DB_URL"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	ProductionMode string `mapstructure:"PROD"`
}

type AwsConfig struct {
	S3Bucket string `mapstructure:"S3_BUCKET"`
	Region   string `mapstructure:"AWS_REGION"`
	Creds1   string `mapstructure:"AWS_CREDENTIALS_1"`
	Creds2   string `mapstructure:"AWS_CREDENTIALS_2"`
	Creds3   string `mapstructure:"AWS_CREDENTIALS_3"`
	ApiKey   string `mapstructure:"API_KEY"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfigApp(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}

func LoadConfigAws(path string, prod string) (config AwsConfig, err error) {
	viper.AddConfigPath(path)
	if prod == "true" {
		viper.SetConfigName("app")
	} else {
		viper.SetConfigName("aws-dev")
	}

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}
