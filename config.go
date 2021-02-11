package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigMap struct {
	SentryDsn          string `yaml:"sentry_dsn"`
	Endpoint           string `yaml:"endpoint"`
	Region             string `yaml:"region"`
	UseHttps           bool   `yaml:"use_https"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	AccessKey          string `yaml:"access_key"`
	SecretKey          string `yaml:"secret_key"`
}

var config = new(ConfigMap)

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("could not decode config file: %v", err)
	}
	log.Printf("ConfigMap Loaded!")
	return nil
}
