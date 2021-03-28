package main

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type ConfigMap struct {
	SentryDsn          string `hcl:"sentry_dsn"`
	Endpoint           string `hcl:"endpoint"`
	Region             string `hcl:"region"`
	UseHttps           bool   `hcl:"use_https"`
	InsecureSkipVerify bool   `hcl:"insecure_skip_verify"`
	AccessKey          string `hcl:"access_key"`
	SecretKey          string `hcl:"secret_key"`
}

var config = new(ConfigMap)

func loadConfig(filename string) (err error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	obj, err := hcl.Parse(string(d))
	if err != nil {
		return err
	}
	// Build up the result
	if err := hcl.DecodeObject(&config, obj); err != nil {
		return err
	}
	return
}
