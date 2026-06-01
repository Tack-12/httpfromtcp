package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	f, err := os.Open("message.txt")

	if err != nil {
		log.Panic("error Occured")
		return
	}

	ch := getLinesChannel(f)

	defer f.Close()

	for i := range ch {
		fmt.Printf("Read: %s \n", i)
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
