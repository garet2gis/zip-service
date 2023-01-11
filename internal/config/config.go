package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RootDirectory    string `env:"ROOT_DIRECTORY" env-default:"root"`
	UploadsDirectory string `env:"UPLOADS_DIRECTORY" env-default:"uploads"`

	HTTPConfig struct {
		Port        string `env:"PORT"  env-required:"true"`
		Host        string `env:"HOST"  env-required:"true"`
		SendTimeout int    `env:"SEND_TIMEOUT" env-default:"0"`
		ReadTimeout int    `env:"READ_TIMEOUT" env-default:"0"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("Read application configuration")

		instance = &Config{}
		if err := cleanenv.ReadConfig(".env", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})

	return instance
}
