package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const port = ":42069"

func main() {

	f, err := net.Listen("tcp", port)

	if err != nil {
		log.Panic("error Occured")
		return
	}

	fmt.Println("Liestening traffic on port ", port)

	for {
		conn, err := f.Accept()

		if err != nil {
			log.Fatal(err)
			break
		}

		fmt.Println("Connection established on ", conn.RemoteAddr())

		ch, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal(err)
			break
		}

		fmt.Println("Request Line:")
		fmt.Printf("Method: %s \n", ch.RequestLine.Method)
		fmt.Printf("Target: %s \n", ch.RequestLine.RequestTarget)
		fmt.Printf("Http Version: %s \n", ch.RequestLine.HttpVersion)

		fmt.Println("Connection ended", conn.RemoteAddr())

	}

}
