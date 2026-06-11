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

		rq, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal(err)
			break
		}

		fmt.Println("Request Line:")
		fmt.Printf("Method: %s \n", rq.RequestLine.Method)
		fmt.Printf("Target: %s \n", rq.RequestLine.RequestTarget)
		fmt.Printf("Http Version: %s \n", rq.RequestLine.HttpVersion)

		fmt.Println("Methods:")

		for key, value := range rq.Headers {
			fmt.Printf("This is called")
			fmt.Printf("- %s: %s \n", key, value)
		}

		fmt.Println("Connection ended", conn.RemoteAddr())

	}

}
