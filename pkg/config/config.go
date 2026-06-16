package config


import (
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/drona-gyawali/runner/pkg/types"
)


func MustLoad() *types.Jobs {
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		log.Fatal("Unable to locate configuration file")
	}

	var cfg types.Jobs
	err:= cleanenv.ReadConfig(config_path, &cfg)

	if err != nil{
		log.Fatalf("Unable to read configuration file %s", err.Error())
	}

	return &cfg
}
