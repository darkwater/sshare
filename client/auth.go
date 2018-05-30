package main

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"os"

	"github.com/darkwater/sshare/common"

	"github.com/golang/protobuf/proto"
)

func getUserDataDir() string {
	// TODO: Properly deduce this path instead of assuming ~/.local/share
	return os.Getenv("HOME") + "/.local/share/sshare"
}

func getPrivateKeyPath() string {
	return getUserDataDir() + "/pkey.bin"
}

// LoadPrivateKey reads the key from disk, or generates a new one if one doesn't exist yet
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	keyFile, err := ioutil.ReadFile(getPrivateKeyPath())
	if err != nil {
		// TODO: Handle errors other than 'file does not exist'
		return CreatePrivateKey()
	}

	keyProto := &common.PrivateKey{}
	if err := proto.Unmarshal(keyFile, keyProto); err != nil {
		return nil, err
	}

	key := keyProto.FromProtobuf()
	return key, nil
}

// CreatePrivateKey creates a new private key and writes it to disk
func CreatePrivateKey() (*rsa.PrivateKey, error) {
	// Generate key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	msg := common.PrivateKey{}
	msg.ToProtobuf(key)
	out, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	// Write to disk
	err = os.MkdirAll(getUserDataDir(), 0700)
	if err != nil {
		return nil, err
	}

	keyFile, err := os.Create(getPrivateKeyPath())
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	keyFile.Write(out)
	keyFile.Chmod(0600)

	return key, nil
}
