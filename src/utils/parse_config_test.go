package utils

import (
	modelConfig "github.com/vladpereskokov/Technopark_HighLoad-nginx/src/models/configs"
	modelServer "github.com/vladpereskokov/Technopark_HighLoad-nginx/src/models/server"
	"strings"
	"testing"
)

const (
	FILE_1 = "../../configs/test_configs/config_test_1.json"
	FILE_2 = "../../configs/test_configs/config_test_2.json"
	FILE_3 = "../../configs/test_configs/config_test_3.json"
	FILE_4 = "../../configs/test_configs/config_test_4.json"
	FILE_5 = "../../configs/test_configs/config_test_5.json"
	FILE_6 = "bladir/bla.json"
	FILE_7 = "../../configs/test_configs/config_test_6.json"
)

func TestFromFileServer(t *testing.T) {
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

	if dirConfig.SetDir() != "/usr/topDir/" {
		t.Errorf("dir don't match! %s", dirConfig.SetDir())
	}
}

func TestFromFileServerPwd(t *testing.T) {
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

	if !strings.Contains(dirConfig.SetDir(), "github.com/vladpereskokov/Technopark_HighLoad-nginx/") {
		t.Errorf("pwd dir don't match! %s", dirConfig.SetDir())
	}
}

func TestFromFileCodes(t *testing.T) {
	config := new(modelServer.Constants)

	err := FromFile(FILE_3, &config)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	for _, value := range *config.Statuses {
		if value.Code != 200 || value.Message != "Ok" {
			t.Error("Don't work")
		}
	}
}

func TestFromFileMethod(t *testing.T) {
	config := new(modelServer.Constants)

	err := FromFile(FILE_4, &config)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	for _, value := range *config.Methods {
		if value.Type != "GET" {
			t.Error("Don't work")
		}
	}
}

func TestFromFileContentType(t *testing.T) {
	config := new(modelServer.Constants)

	err := FromFile(FILE_5, &config)
	if err != nil {
		t.Errorf("%v\n", err.Error())
	}

	for _, value := range *config.ContentTypes {
		if value.Extension != ".css" || value.Type != "text/css" {
			t.Error("Don't work")
		}
	}
}

func TestFromFileFail(t *testing.T) {
	config := new(modelServer.Constants)

	err := FromFile(FILE_6, &config)
	if err == nil {
		t.Error("Dir exist!")
	}
}

func TestFromFileErrorFile(t *testing.T) {
	config := new(modelServer.Constants)

	err := FromFile(FILE_7, &config)
	if err == nil {
		t.Error("Json cool")
	}
}
