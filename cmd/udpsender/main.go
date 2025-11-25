package main

import (
	"bufio"
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
		res, err := reader.ReadString('\n')
		logError(err)
		conn.Write([]byte(res))
	}
}
