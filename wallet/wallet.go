package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/PulseCoinOrg/nexacoin/common"
)

type Wallet struct {
	PublicKey  []byte
	PrivateKey []byte
	Address    common.Address
}

func New() (*Wallet, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privKeyBytes := key.D.Bytes()

	pubKey := key.PublicKey
	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)

	return &Wallet{
		PublicKey:  pubKeyBytes,
		PrivateKey: privKeyBytes,
		Address:    common.MakeAddr(pubKeyBytes),
	}, nil
}

func (w *Wallet) PublicKeyBytes() []byte {
	return w.PublicKey[:]
}
