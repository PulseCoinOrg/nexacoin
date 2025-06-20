package core

import (
	"errors"
	"time"

	"github.com/PulseCoinOrg/nexacoin/core/types"
	"github.com/PulseCoinOrg/nexacoin/wallet"
)

type Validator struct {
	Wallet          *wallet.Wallet
	CurrentBlock    *types.Block
	ValidatedBlocks []*types.Block
}

func NewValidator() (*Validator, error) {
	var validated []*types.Block
	wallet, err := wallet.LoadFromDisk()
	if err != nil {
		return nil, err
	}
	return &Validator{
		Wallet:          wallet,
		ValidatedBlocks: validated,
	}, nil
}

func (v *Validator) GetValidatorAddress() (string, error) {
	if v.Wallet == nil {
		return "", errors.New("failed to get validators address as wallet is nil")
	}
	return v.Wallet.Address.Hex(), nil
}

func (v *Validator) ValidateBlock(b *types.Block) bool {
	if b == nil {
		return false
	}

	now := time.Now().Unix()
	if b.Time > now+10*int64(time.Minute) || b.Time < 0 {
		return false
	}

	for _, vb := range v.ValidatedBlocks {
		if vb.Height == b.Height && vb.Hash.Hex() != b.Hash.Hex() {
			return false
		}
	}

	//for _, tx := range b.Transactions {
	// 		if !tx.IsValid() {
	// 			return false
	//		}
	//}

	v.ValidatedBlocks = append(v.ValidatedBlocks, b)

	return true
}
