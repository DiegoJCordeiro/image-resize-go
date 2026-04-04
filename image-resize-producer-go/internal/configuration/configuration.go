package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Configuration struct {
	SRVRConfiguration struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	AWSConfig struct {
		Endpoint string `mapstructure:"endpoint"`
	} `mapstructure:"aws"`
	DBConfiguration struct {
		RedisConfiguration struct {
			Host     string `mapstructure:"host"`
			Password string `mapstructure:"password"`
			DBN      int    `mapstructure:"dbn"`
		} `mapstructure:"redis"`
		MongoDBConfiguration struct {
			Host     string `mapstructure:"host"`
			DBName   string `mapstructure:"name"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"mongodb"`
	} `mapstructure:"database"`
	KafkaConfiguration struct {
		BootstrapServer string `mapstructure:"bootstrap-server"`
		KafkaAcks       struct {
			ProducerAck string `mapstructure:"producer"`
			ConsumerAck string `mapstructure:"consumer"`
		} `mapstructure:"acknowledge"`
		KafkaTopics struct {
			ImageProcess struct {
				Name string `mapstructure:"name"`
			} `mapstructure:"image-process"`
		} `mapstructure:"topics"`
	} `mapstructure:"apache-kafka"`
}

func LoadConfiguration(profile string) (*Configuration, error) {

	if profile != "" {
		if err := godotenv.Load(".env." + profile); err != nil {
			return nil, err
		}

		v := viper.New()
		v.SetConfigName("application-" + profile)
		v.SetConfigType("yaml")
		v.AddConfigPath(".")

		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}

		for _, key := range v.AllKeys() {
			val := v.GetString(key)
			if strings.Contains(val, "${") {
				v.Set(key, os.ExpandEnv(val))
			}
		}

		var cfg Configuration

		if err := v.Unmarshal(&cfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return &cfg, nil
	} else {
		return nil, fmt.Errorf("invalid configuration profile")
	}
}
