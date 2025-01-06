package util

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	SecretKey      string `mapstructure:"SECRET_KEY"`
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
	mode := os.Getenv("MODE")
	viper.AddConfigPath(path)

	if mode == "test" {
		fmt.Println("test mode")
		viper.SetConfigName("test")
		viper.SetConfigType("env")
	} else {
		viper.SetConfigName("app")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	fmt.Println(config.DBUrl)
	return
}
