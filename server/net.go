package main

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/darkwater/sshare/common"
	"github.com/golang/protobuf/proto"
)

// Actions a client can take at the start of a connection
const (
	welcomeActionAuthenticate byte = iota
	welcomeActionUseInvite
)

// startListening starts a TLS server and returns its handle
func startListening() (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{Certificates: []tls.Certificate{cert}}
	sock, err := tls.Listen("tcp", "127.0.0.1:3636", conf)
	if err != nil {
		return nil, err
	}

	return sock, nil
}

func welcome(conn net.Conn) (int, error) {
	welcome := &common.Welcome{
		Url:             "s.dark.red",
		AcceptsNewUsers: true,
		AuthChallenge:   make([]byte, 16),
	}

	// Create an auth challenge
	rand.Read(welcome.AuthChallenge)

	// Send Welcome
	common.SendMessage(conn, common.MsgWelcome, welcome)

	// Read response
	msgtype, msg, err := common.ReadMessage(conn)
	if err != nil {
		return 0, err
	}

	switch msgtype {
	case common.MsgAuthResponse:
		// Attempt authentication
		response := &common.AuthResponse{}
		if err := proto.Unmarshal(msg, response); err != nil {
			return 0, err
		}

		return authenticateClient(response, conn)

	case common.MsgInviteUse:
		// Attempt registration
		response := &common.InviteUse{}
		if err := proto.Unmarshal(msg, response); err != nil {
			return 0, err
		}

		return useInvite(response, conn)

	default:
		// Invalid response
		return 0, errors.New("invalid response to welcome")
	}
}

func authenticateClient(msg *common.AuthResponse, conn net.Conn) (int, error) {
	// TODO: This algorithm is similar to SSHv1, but SSHv2 uses a different
	// algorithm - why?

	fmt.Printf("Response: %x\n", msg.Signature)

	return 0, nil
}

func useInvite(msg *common.InviteUse, conn net.Conn) (int, error) {
	invite, err := dbGetInvite(msg.Code)
	if err != nil {
		return 0, errors.New("invalid invitation code")
	}

	user := &User{
		URL:       "xf",
		InvitedBy: invite.Sender,
	}

	dbUseInvite(invite, user)
	dbAddPublicKey(user, msg.Key)

	return 0, nil
}
