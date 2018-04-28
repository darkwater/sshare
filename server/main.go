package main

import (
	"crypto/tls"
	"io/ioutil"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}

	conf := &tls.Config{Certificates: []tls.Certificate{cert}}
	sock, err := tls.Listen("tcp", "127.0.0.1:3636", conf)
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
	out, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}

	println(string(out))
	conn.Write(out)
	conn.Close()
}
