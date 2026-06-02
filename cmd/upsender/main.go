package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	addr, err := net.ResolveUDPAddr("udp", ":42069")

	if err != nil {
		fmt.Printf("There was an error Resolving UDP Address:%v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Printf("There was an error with the Dial up:%v", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")

		value, err := reader.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				break
			}
			fmt.Printf("There was an error Reading the string value from the console:%v", err)
		}

		_, err = conn.Write([]byte(value))

		if err != nil {
			fmt.Printf("There was an error Writing to the UDP connection :%v", err)
		}

	}

}
