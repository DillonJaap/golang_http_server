package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func noError(err error, messages ...string) {
	if err != nil {
		panic(append(messages, err.Error()))
	}
}
func logError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	noError(err)

	conn, err := net.DialUDP("udp", nil, addr)
	noError(err)
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		res, err := reader.ReadString('\n')
		logError(err)
		conn.Write([]byte(res))
	}
}
