package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

//Configuration - struct holds the config throughout the application
type Configuration struct {
	AppName  string `json:"app_name"`
	Database struct {
		Type      string  `json:"type"`
		Name      string  `json:"name"`
		URI       string  `json:"uri"`
		Timeout   int     `json:"timeout"`
		PoolLimit *uint64 `json:"pool_limit"`
		Source    string  `json:"source"`
		UserName  string  `json:"user_name"`
		Password  string  `json:"password"`
	} `json:"database"`
	RedisVars struct {
		ConnString      string `json:"conn_string"`
		MaxRetries      int    `json:"max_retries"`
		MinRetryBackoff int    `json:"min_retry_backoff"`
		DialTimeout     int    `json:"dial_timeout"`
		ReadTimeout     int    `json:"read_timeout"`
		WriteTimeout    int    `json:"write_timeout"`
		Password        string `json:"password"`
	} `json:"redis_vars"`
}

var (
	Cfg  *Configuration
	once sync.Once
)

// Parse parses the json configuration file
// And converting it into native type
func Parse(file string) *Configuration {
	once.Do(func() {
		// Reading the flags
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("config: ioutil.ReadFile failed: %s", err.Error())
		}

		if err := json.Unmarshal(data, &Cfg); err != nil {
			log.Fatalf("config: json.unmarshal failed: %s", err.Error())
		}
	})

	return Cfg
}
