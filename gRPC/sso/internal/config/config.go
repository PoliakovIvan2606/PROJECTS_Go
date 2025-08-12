package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 		string 			`yaml:"env" env-default:"local"`
	StoragePath string 			`yaml:"storage_path" env-required:"true"`
	TokenTTL	time.Duration 	`yaml:"token_ttl" env-required:"1h"`
	GRPC		GRPCConfig 		`yaml:"grpc"`
}

type GRPCConfig struct {
	Port	int 			`yaml:"port"`
	Timwout time.Duration 	`yaml:"timeout"`
}

func MastLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MastLoadByPath(path)
}

func MastLoadByPath(path string) *Config {
	// проверка на то что есть файл
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exost: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("field to read config: " + err.Error())
	}

	return &cfg
}



// Priority: flug > env > default
// flug: --config=./path...
// env: CONFIG_PATH=./path/...
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}


