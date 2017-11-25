package common

import (
	"encoding/json"
	"os"
)

type Config struct {
	MySQLUser     string `json:"mysql_user"`
	MySQLPass     string `json:"mysql_pass"`
	MySQLDatabase string `json:"mysql_database"`
}

func ReadConfig(path string) (*Config, error) {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	c := &Config{}
	err := decoder.Decode(c)
	return c, err
}
