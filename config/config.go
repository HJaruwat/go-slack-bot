package config

import (
	"fmt"
	"log"

	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

type (
	//Conf struct of
	Conf struct {
		App struct {
			Queue Queue `mapstructure:"queue"`
		}
	}

	//Queue struct of
	Queue struct {
		Host     string `mapstructure:"host"`
		Tool     string `mapstructure:"tool"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
)

// Setting is instance config setting
var Setting Conf

var stage string

var opts Options

// Options is
type Options struct {
	Environment string `short:"e" long:"env" description:"Environment for load config" required:"true"`
}

func init() {
	_, err := flags.Parse(&opts)

	if err != nil {
		log.Fatal(err)
		return
	}

	stage = ""
	if opts.Environment == "dev" || opts.Environment == "development" {
		stage = "development"
	} else {
		stage = "production"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf(".env.%s", stage))
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Setting)
	if err != nil {
		panic(err)
	}

	log.Printf(" [*] Env name : %s", stage)
}
