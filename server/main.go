package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/darkwater/sshare/common"
	proto "github.com/golang/protobuf/proto"
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
	// TODO: This algorithm is similar to SSHv1, but SSHv2 uses a different
	// algorithm - why?
	challenge := &common.AuthChallenge{
		Nonce: make([]byte, 16),
	}

	// Generate a random 32-byte challenge
	rand.Read(challenge.Nonce)

	common.SendMessage(conn, common.MsgAuthChallenge, challenge)

	// Read response
	msgtype, msg, err := common.ReadMessage(conn)
	if err != nil {
		panic(err)
	}

	// Assert that the message we got is a response
	if msgtype != common.MsgAuthResponse {
		panic("unexpected msgtype " + string(msgtype))
	}

	response := &common.AuthResponse{}
	if err := proto.Unmarshal(msg, response); err != nil {
		panic(err)
	}

	fmt.Printf("Response: %x\n", response.Signature)

	conn.Close()
}
