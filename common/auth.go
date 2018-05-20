package common

import (
	"crypto/rsa"
	"crypto/sha256"
	"math/big"
	"strconv"
)

// ToProtobuf populates the message's fields from a crypto/rsa PrivateKey
func (p *PrivateKey) ToProtobuf(key *rsa.PrivateKey) {
	p.D = key.D.Text(62)
	p.Primes = make([]string, len(key.Primes))
	for k, prime := range key.Primes {
		p.Primes[k] = prime.Text(62)
	}

	p.Public = &PublicKey{}
	p.Public.ToProtobuf(&key.PublicKey)
}

// FromProtobuf creates a crypto/rsa PrivateKey from this protobuf message
func (p *PrivateKey) FromProtobuf() *rsa.PrivateKey {
	key := rsa.PrivateKey{
		PublicKey: *p.Public.FromProtobuf(),
		D:         &big.Int{},
		Primes:    make([]*big.Int, len(p.Primes)),
	}

	key.D.SetString(p.D, 62)
	for k, prime := range p.Primes {
		key.Primes[k] = &big.Int{}
		key.Primes[k].SetString(prime, 62)
	}

	return &key
}

// ToProtobuf populates the message's fields from a crypto/rsa PublicKey
func (p *PublicKey) ToProtobuf(key *rsa.PublicKey) {
	p.N = key.N.Text(62)
	p.E = int64(key.E)
	// p.hash = HashPublicKey(key)
}

// FromProtobuf creates a crypto/rsa PublicKey from this protobuf message
func (p *PublicKey) FromProtobuf() *rsa.PublicKey {
	pubkey := rsa.PublicKey{
		N: &big.Int{},
		E: int(p.E),
	}

	pubkey.N.SetString(p.N, 62)

	return &pubkey
}

// HashPublicKey uses SHA-256 to create a hash of the public key
func HashPublicKey(key *rsa.PublicKey) [sha256.Size]byte {
	str := key.N.Text(62) + "-" + strconv.Itoa(key.E)
	hash := sha256.Sum256([]byte(str))
	return hash
}
