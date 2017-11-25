package common

import (
	"encoding/json"
	"os"
)

type Config struct {
	MySQLHost     string `json:"mysql_host"`
	MySQLPort     string `json:"mysql_port"`
	MySQLDatabase string `json:"mysql_database"`
	MySQLUser     string `json:"mysql_user"`
	MySQLPass     string `json:"mysql_pass"`
}

func ReadConfig(path string) (*Config, error) {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	c := &Config{}
	err := decoder.Decode(c)
	return c, err
}
