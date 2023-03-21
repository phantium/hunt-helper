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
	Notifications struct {
		SlackAPIKey    string `yaml:"slack_api_key"`
		SlackChannelID string `yaml:"slack_channel_id"`
	} `yaml:"notifications"`
	WebServer struct {
		Host               string `yaml:"host"`
		Port               string `yaml:"port"`
		CertificateFile    string `yaml:"certificate_file"`
		CertificateKeyFile string `yaml:"certificate_key_file"`
		StaticDir          string `yaml:"static_dir"`
		TemplatesDir       string `yaml:"templates_dir"`
	} `yaml:"webserver"`
	GCP struct {
		GCPCredentialsFile string `yaml:"gcp_credentials_file"`
	} `yaml:"gcp"`
	Storage struct {
		ResticGCSConfig string `yaml:"restic_gcs_config"`
		ResticGCSBucket string `yaml:"restic_gcs_bucket"`
	} `yaml:"storage"`
	Psono struct {
		PsonoConfig   string `yaml:"psono_config"`
		PsonoSecretId string `yaml:"psono_secret_id"`
	} `yaml:"psono"`
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
