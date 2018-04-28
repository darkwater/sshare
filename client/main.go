package main

import (
	"crypto/tls"
	"io/ioutil"
)

func main() {
	certFile, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		panic(err)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{certFile},
	}
	_ = cert

	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	sock, err := tls.Dial("tcp", "127.0.0.1:3636", &conf)
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	sock.Write([]byte("hello world"))
	sock.CloseWrite()

	read, err := ioutil.ReadAll(sock)
	println(string(read), err)

	sock.Close()
}
