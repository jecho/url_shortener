package main

import (
	"testing"
)

func TestNewDB(t *testing.T) {
	configFile := ".env/foxley_mock.json"
	config := loadConfig(configFile)

	dsn := config.DB_USER + ":" + config.DB_PASS + "@" + config.DB_HOST + "/" + config.DB_NAME
	_, err := NewDB(dsn)

	if err != nil {
		t.Errorf("Could not establish connection to mock database.")
	}
}
