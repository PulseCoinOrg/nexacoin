/*
 * NexaCoin - A Cryptocurrency Framework
 *
 * Copyright (c) 2025 NexaCoin Developers
 *
 * This file is part of the NexaCoin project and is licensed under the MIT License.
 * You may obtain a copy of the License at:
 *
 *     https://opensource.org/licenses/MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"math/big"

	"github.com/PulseCoinOrg/nexacoin/common"
)

var (
	DefaultWalletPath = "./wallet.key"
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

func (w *Wallet) SaveDisk() error {
	privKeyHex := hex.EncodeToString(w.PrivateKey)
	return ioutil.WriteFile(DefaultWalletPath, []byte(privKeyHex), 0600)
}

func LoadFromDisk() (*Wallet, error) {
	data, err := ioutil.ReadFile(DefaultWalletPath)
	if err != nil {
		return nil, errors.New("wallet not found; run `wallet.New()` and save it first")
	}

	privKeyBytes, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, errors.New("failed to decode private key")
	}

	priv := new(ecdsa.PrivateKey)
	priv.D = new(big.Int).SetBytes(privKeyBytes)
	priv.PublicKey.Curve = elliptic.P256()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(privKeyBytes)

	pubKeyBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)

	return &Wallet{
		PrivateKey: privKeyBytes,
		PublicKey:  pubKeyBytes,
		Address:    common.MakeAddr(pubKeyBytes),
	}, nil
}

func (w *Wallet) PublicKeyBytes() []byte {
	return w.PublicKey[:]
}
