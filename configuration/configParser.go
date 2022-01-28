package configuration

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	Images []Image     `yaml:"images"`
	Email  EmailConfig `yaml:"email"`
}

type EmailConfig struct {
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Image struct {
	Name        string `yaml:"name"`
	Constraint  string `yaml:"constraint"`
	UsedVersion string `yaml:"usedVersion"`
}

func ParseConfig(file string) (*Configuration, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	configuration := new(Configuration)
	err = yaml.Unmarshal(data, configuration)

	return configuration, err
}
