package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

var ServerClosed atomic.Bool

type Server struct {
	Listner net.Listener
}

func Serve(port int) (*Server, error) {

	var MainServer *Server

	portNum := strconv.Itoa(port)
	server, err := net.Listen("tcp", ":"+portNum)

	if err != nil {
		return nil, fmt.Errorf("There was an error creating the server in port %v: Error: %s", port, err)
	}

	MainServer = &Server{
		Listner: server,
	}

	go MainServer.listenConnection()

	return MainServer, nil
}

func (s *Server) listenConnection() {

	for {
		if !ServerClosed.Load() {
			conn, err := s.Listner.Accept()

			if err != nil {
				log.Fatalf("Error getting a connection. Error: %s", err)
			}

			go s.handleConn(conn)
		}

	}
}

func (s *Server) Close() error {
	err := s.Listner.Close()
	ServerClosed.Store(true)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) handleConn(conn net.Conn) {

	response := []byte("HTTP/1.1 200 OK\nContent-Type: text/plain\nContent-Length: 13\nHello World!\n")

	_, err := conn.Write(response)

	if err != nil {
		log.Fatalf("Error Writing into the Connection %s", err)
	}

}
