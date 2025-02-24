package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	CLOUD_FLARE_ACCESS_KEY_ID string `mapstructure:"CLOUD_FLARE_ACCESS_KEY_ID"`
	CLOUD_FLARE_SECRET_ACCESS_KEY string `mapstructure:"CLOUD_FLARE_SECRET_ACCESS_KEY"`
	CLOUD_FLARE_TOKEN string `mapstructure:"CLOUD_FLARE_TOKEN"`
	CLOUD_FLARE_R2_BUCKET string `mapstructure:"CLOUD_FLARE_R2_BUCKET"`
	CLOUD_FLARE_R2_REGION string `mapstructure:"CLOUD_FLARE_R2_REGION"`
	CLOUD_FLARE_R2_ENDPOINT string `mapstructure:"CLOUD_FLARE_R2_ENDPOINT"`
	BACKUP_DATABASE_URL string `mapstructure:"BACKUP_DATABASE_URL"`
	BACKUP_DATABASE_PASSWORD string `mapstructure:"BACKUP_DATABASE_PASSWORD"`
	BACKUP_DATABASE_USER string `mapstructure:"BACKUP_DATABASE_USER"`
	BACKUP_DATABASE_HOST string `mapstructure:"BACKUP_DATABASE_HOST"`
	BACKUP_DATABASE_PORT string `mapstructure:"BACKUP_DATABASE_PORT"`
	BACKUP_DATABASE_NAME string `mapstructure:"BACKUP_DATABASE_NAME"`
	BACKUP_CRON_SCHEDULE string `mapstructure:"BACKUP_CRON_SCHEDULE"`
	RUN_ON_STARTUP bool `mapstructure:"RUN_ON_STARTUP"`
	SINGLE_SHOT_MODE bool `mapstructure:"SINGLE_SHOT_MODE"`
	BACKUP_FILE_PREFIX string `mapstructure:"BACKUP_FILE_PREFIX"`
	BUCKET_SUBFOLDER string `mapstructure:"BUCKET_SUBFOLDER"`
}

// Loads app configuration from .env file.
func LoadConfig(path ,name, fileType string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(fileType)
	setDefaults()

	viper.AutomaticEnv()

	// if file not found use authomatic env
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} else {
			return Config{}, err
		}
	}

	var config Config

	return config, viper.Unmarshal(&config)
}

func setDefaults() {
	viper.SetDefault("CLOUD_FLARE_ACCESS_KEY_ID", "")
	viper.SetDefault("CLOUD_FLARE_SECRET_ACCESS_KEY", "")
	viper.SetDefault("CLOUD_FLARE_R2_BUCKET", "")
	viper.SetDefault("CLOUD_FLARE_R2_REGION", "")
}