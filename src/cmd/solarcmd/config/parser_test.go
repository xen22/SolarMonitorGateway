package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func writeConfig(filename string, contents string) {
	ioutil.WriteFile(filename, []byte(contents), 0644)
}

func TestParseNoConfigFile(t *testing.T) {
	_, err := ParseConfig("non_existent_file")
	if err == nil {
		t.Errorf("Expected test to return error 'file not found', but it didn't.")
	}
}

func TestParseEmptyConfigFile(t *testing.T) {
	f := "/tmp/test.json"
	writeConfig(f, "")
	defer os.Remove(f)

	_, err := ParseConfig(f)
	if err == nil {
		t.Errorf("Expected test to return error 'file not found', but it didn't.")
	}
}

// func TestParseSimple(t *testing.T) {
// 	var f = "test.json"
// 	_, err := ParseConfig(f)
// 	if err != nil {
// 		t.Errorf("Parsing failed, err: %s", err)
// 	}
// }

func TestParseDb(t *testing.T) {
	f := "/tmp/test.json"
	confdata :=
		`{
		"Database": {
			"DbName": "SolarMonitorDb",
			"DbUser": "u1",
			"DbPassword": "ppp",
			"SiteName": "Test system"
		},
		"Sensors": []
	}
	`
	writeConfig(f, confdata)
	defer os.Remove(f)

	c, err := ParseConfig(f)
	if err != nil {
		t.Errorf("Parsing failed, err: %s", err)
	}

	if c.Database.Name != "SolarMonitorDb" {
		t.Errorf("Db name incorrect")
	}
	if c.Database.User != "u1" {
		t.Errorf("Db user incorrect")
	}
	if c.Database.Password != "ppp" {
		t.Errorf("Db password incorrect")
	}
	if c.Database.SiteName != "Test system" {
		t.Errorf("SiteName incorrect")
	}
}
