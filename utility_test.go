package main

import (
	"testing"
	"reflect"
)

func TestLoadConfig(t *testing.T) {

	configFile := ".env/foxley_mock.json"
	configTest := loadConfig(configFile)

	configMock := Config{
		DB_HOST: "tcp(foxley.co:22113)",
		DB_NAME: "spellweaver",
		DB_USER: "root",
		DB_PASS: "zelda",
	}

	got := reflect.DeepEqual(configMock, configTest)
	if got != true {
		t.Errorf("configMock does not equal configTest")
	}
}
