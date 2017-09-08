package handler

import (
	"fmt"
	modelServer "github.com/vladpereskokov/Technopark_HighLoad-nginx/src/models/server"
	"github.com/vladpereskokov/Technopark_HighLoad-nginx/src/utils"
	"log"
	"net"
	"net/url"
)

const HTTP_CONSTANTS_CONFIG = "configs/http.json"

type Handler struct {
	Connection net.Conn
	Request    *modelServer.Request
	Response   *modelServer.Response
	Constants  *modelServer.Constants
	Dir        string
}

type HandlerFunc func(chan net.Conn)

func CreateHandler(dir string) (handlerFunc HandlerFunc) {
	httpConstantsConfig := new(modelServer.Constants)

	err := utils.FromFile(HTTP_CONSTANTS_CONFIG, &httpConstantsConfig)
	if err != nil {
		log.Panicf("can not init http config: %v", err)
	}

	return func(channel chan net.Conn) {
		handler := Handler{}
		handler.create(httpConstantsConfig, dir)

		go handler.start(channel)
	}
}

func (handler *Handler) create(config *modelServer.Constants, dir string) {
	handler.Request = modelServer.InitRequest()

	handler.Response = new(modelServer.Response)
	handler.Response.Status = new(modelServer.Status)
	handler.Response.Headers = make(map[string]string)

	handler.Response.Status.Message = "ok"
	handler.Response.Status.Code = 200

	handler.Constants = config
	handler.Dir = dir
}

func (handler *Handler) start(channel chan net.Conn) {
	for {
		handler.Connection = <-channel
		handler.handle()
	}
}

func (handler *Handler) handle() {
	handler.readRequest()
	handler.closeConn()
}

func (handler *Handler) readRequest() {
	buffer := make([]byte, 1024)

	_, err := handler.Connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())

		return
	}

	handler.Connection.Write(buffer)
	handler.Connection.Write([]byte("\r\n\r\n"))
}

func (handler *Handler) closeConn() {
	handler.Connection.Close()
}
