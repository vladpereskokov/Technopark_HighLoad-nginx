package server

import (
	"fmt"
	"github.com/vladpereskokov/Technopark_HighLoad-nginx/src/handler"
	modelConfig "github.com/vladpereskokov/Technopark_HighLoad-nginx/src/models/configs"
	"log"
	"net"
	"os"
)

type Server struct {
	network  string
	protocol string
	host     string
	port     string
	isSetup  bool
}

func (server *Server) CreateServer(config modelConfig.Server) {
	server.setNetwork(config.Network)
	server.setProtocol(config.Protocol)
	server.setHost(config.Host)
	server.setPort(string(config.Port))

}

func (server *Server) Start(config *modelConfig.Config) {
	serverConf := config.GetServer()

	listener, err := net.Listen(serverConf.GetNetwork(), ":"+serverConf.GetPort())
	if err != nil {
		panic("Failed start server: " + err.Error())
	}

	defer listener.Close()

	log.Print("Server started at " + serverConf.GetPort() + " port")

	ch := make(chan net.Conn)

	handle := handler.Handler{}

	for i := 0; i < 4; i++ {
		go handle.Start(ch)
		println("Created worker...")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		ch <- conn
	}
}

func (server *Server) setNetwork(network string) {
	server.network = network
}

func (server *Server) setProtocol(protocol string) {
	server.protocol = protocol
}

func (server *Server) setHost(host string) {
	server.host = host
}

func (server *Server) setPort(port string) {
	server.port = port
}

func (server *Server) setSetup(isSetup bool) {
	server.isSetup = isSetup
}
