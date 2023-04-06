package configuration

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DiscordConfig struct {
	Discord struct {
		Token string `yaml:"token"`
	} `yaml:"discord"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Timeout  string `yaml:"timeout"`
	} `yaml:"database"`
}

// Read file and decode yaml using the provided interface
func ReadConfig(cfg interface{}, cfg_file string) {
	f, err := os.Open(cfg_file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		fmt.Println(err)
	}
}

// Read file and return the file as a string
func ReadFile(filename string) string {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(file)
}
