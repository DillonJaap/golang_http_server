package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func assert(condition bool, messages ...string) {
	if !condition {
		panic(messages)
	}
}

func noError(err error, messages ...string) {
	if err != nil {
		panic(append(messages, err.Error()))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	noError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		noError(err)

		fmt.Println("A connection has been accepted")
		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("The connection has been closed")
		conn.Close()
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		buffer := make([]byte, 8)
		current := ""
		for {
			n, err := f.Read(buffer)
			if errors.Is(err, io.EOF) {
				break
			}
			noError(err, "failed to read file")

			parts := strings.Split(string(buffer[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				out <- current + parts[i]
				current = ""
			}
			current += parts[len(parts)-1]
		}
	}()
	return out
}
