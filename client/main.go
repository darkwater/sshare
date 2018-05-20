package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/darkwater/sshare/common"

	"github.com/golang/protobuf/proto"
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

	conn, err := tls.Dial("tcp", "127.0.0.1:3636", &conf)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	handle(conn)
}

func handle(conn net.Conn) {
	// Read challenge
	msgtype, msg, err := common.ReadMessage(conn)
	if err != nil {
		panic(err)
	}

	// Assert that the message we got is a challenge
	if msgtype != common.MsgWelcome {
		panic("unexpected msgtype " + string(msgtype))
	}

	welcome := &common.Welcome{}
	if err := proto.Unmarshal(msg, welcome); err != nil {
		panic(err)
	}

	fmt.Printf("Challenge nonce: %x\n", welcome.AuthChallenge)

	// Sign nonce
	key, err := LoadPrivateKey()
	if err != nil {
		panic(err)
	}

	hash := common.HashPublicKey(&key.PublicKey)
	fmt.Printf("hash: %x\n", hash)

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, 0, welcome.AuthChallenge)
	if err != nil {
		panic(err)
	}

	// Create response
	response := &common.AuthResponse{
		Signature: signature,
	}

	common.SendMessage(conn, common.MsgAuthResponse, response)
}
