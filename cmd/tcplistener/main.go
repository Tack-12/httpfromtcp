package main

import (
	"errors"
	"fmt"
	"io"
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

	fmt.Printf("Liestening traffic on port %s ", port)

	for {
		conn, err := f.Accept()

		if err != nil {
			log.Fatal(err)
			break
		}

		fmt.Printf("Connection established on %s ", conn.RemoteAddr())

		ch := getLinesChannel(conn)

		for i := range ch {
			fmt.Println(i)
		}

	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	buffer := make([]byte, 8)
	ch := make(chan string)

	var line string

	go func() {
		for {

			n, err := f.Read(buffer)

			if err != nil {
				if errors.Is(err, io.EOF) {
					close(ch)
					return
				}
				fmt.Println("Error reading content")
				break
			}

			//Here the n = Length of buffer (Can be filled and be 8 or the amount left)
			for _, cb := range buffer[:n] {
				if cb == '\n' {
					ch <- line
					line = ""
				} else {
					line += string(cb)
				}
			}

		}
	}()
	return ch
}
