package main

import (
	"fmt"
	"net"
)

func main() {
	err := initDatabase()
	if err != nil {
		panic(err)
	}

	sock, err := startListening()
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			panic(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	user, err := welcome(conn)
	if err != nil {
		fmt.Printf("auth failed: %s\n", err)
		conn.Close()
		return
	}

	fmt.Printf("authenticated: %+v\n", user)

	conn.Close()
}
