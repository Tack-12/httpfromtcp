package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
)

type HandlerError struct {
	ErrorCode    response.StatusCode
	ErrorMessage error
}

type Handler func(io.Writer, *request.Request) *HandlerError

type Server struct {
	Listner net.Listener
	Status  bool
	Handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	var MainServer *Server

	portNum := strconv.Itoa(port)
	server, err := net.Listen("tcp", ":"+portNum)

	if err != nil {
		return nil, fmt.Errorf("There was an error creating the server in port %v: Error: %s", port, err)
	}

	MainServer = &Server{
		Listner: server,
		Status:  true,
	}

	go MainServer.listenConnection()

	return MainServer, nil
}

func (s *Server) listenConnection() {

	for {

		conn, err := s.Listner.Accept()
		if !s.Status {
			return
		}
		if err != nil {
			log.Fatalf("Error getting a connection. Error: %s", err)
		}

		go s.handleConn(conn)

	}
}

func (s *Server) Close() error {

	log.Printf("THE CONNECTION is CLOSED")

	var err error

	defer func() {
		err = s.Listner.Close()
	}()

	if err != nil {
		return err
	}
	s.Status = false

	return nil
}

func (s *Server) handleConn(conn net.Conn) {

	if s.Status {

		err := response.WriteStatusLine(conn, 200)

		if err != nil {
			s.Handler{}
		}
		headers := response.GetDefaultHeaders(0)

		err = response.WriteHeaders(conn, headers)

		if err != nil {
			log.Fatalf("Error Writing Headers into the Connection %s", err)
		}
		defer conn.Close()

	} else {
		log.Fatalf("The connection was closed before writing for some reason")
	}

}

func HandleError(io.Writer, *request.Request) *HandlerError {

}
