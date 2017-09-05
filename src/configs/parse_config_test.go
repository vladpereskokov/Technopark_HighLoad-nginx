package configs

import (
	"fmt"
	modelConfig "github.com/vladpereskokov/Technopark_HighLoad-nginx/src/models/configs"
	"strings"
	"testing"
)

const (
	FILE_1 = "../../test_configs/config_test_1.json"
	FILE_2 = "../../test_configs/config_test_2.json"
	FILE_3 = "../../test_configs/config_test_3.json"
)

func TestFromTop(t *testing.T) {
	config := new(modelConfig.Statuses)

	err := FromFile(FILE_3, &config)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	fmt.Printf("\n")
	fmt.Printf("%v", config.Status[0].Code)
	t.Error("hey")
	fmt.Printf("\n")
	fmt.Printf("\n")
}

func TestFromFileFirst(t *testing.T) {
	config := new(modelConfig.Config)

	err := FromFile(FILE_1, &config)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	serverConfig := config.GetServer()
	dirConfig := config.GetDir()

	if serverConfig.GetNetwork() != "tcp" {
		t.Errorf("network don't match! %s", serverConfig.GetNetwork())
	}

	if serverConfig.GetProtocol() != "http" {
		t.Errorf("protocol don't match! %s", serverConfig.GetProtocol())
	}

	if serverConfig.GetHost() != "localhost" {
		t.Errorf("host don't match! %s", serverConfig.GetHost())
	}

	if serverConfig.GetPort() != "2007" {
		t.Errorf("port don't match! %s", serverConfig.GetPort())
	}

	if dirConfig.GetDir() != "/usr/topDir/" {
		t.Errorf("dir don't match! %s", dirConfig.GetDir())
	}
}

func TestFromFileSecond(t *testing.T) {
	config := new(modelConfig.Config)

	err := FromFile(FILE_2, &config)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	serverConfig := config.GetServer()
	dirConfig := config.GetDir()

	if serverConfig.GetNetwork() != "tcp" {
		t.Errorf("network don't match! %s", serverConfig.GetNetwork())
	}

	if serverConfig.GetProtocol() != "http" {
		t.Errorf("protocol don't match! %s", serverConfig.GetProtocol())
	}

	if serverConfig.GetHost() != "localhost" {
		t.Errorf("host don't match! %s", serverConfig.GetHost())
	}

	if serverConfig.GetPort() != "2007" {
		t.Errorf("port don't match! %s", serverConfig.GetPort())
	}

	if !strings.Contains(dirConfig.GetDir(), "github.com/vladpereskokov/Technopark_HighLoad-nginx/") {
		t.Errorf("pwd dir don't match! %s", dirConfig.GetDir())
	}
}
