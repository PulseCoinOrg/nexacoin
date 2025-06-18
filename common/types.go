package common

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	HashLength    = 32
	AddressLength = 20
)

type Hash [HashLength]byte

func SHA256(data []byte) Hash {
	return sha256.Sum256(data)
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) Base58Encode() string {
	return base58.Encode(h[:])
}

type Address [AddressLength]byte

func MakeAddr(pubKeyBytes []byte) Address {
	hasher := ripemd160.New()
	hasher.Write(pubKeyBytes)
	return Address(hasher.Sum(nil)[:AddressLength])
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Hex() string {
	return hex.EncodeToString(a[:])
}
