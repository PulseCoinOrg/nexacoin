package core

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/PulseCoinOrg/nexacoin/common"
	"github.com/PulseCoinOrg/nexacoin/nexadb/memorydb"
)

var (
	ErrNoValidators = errors.New("no validators")
)

type ValidatorPool struct {
	Database          *memorydb.Database
	Validators        map[string]*Validator
	SelectedValidator *Validator
}

func NewValidatorPool() *ValidatorPool {
	db := memorydb.New()
	validators := make(map[string]*Validator)
	return &ValidatorPool{
		Database:   db,
		Validators: validators,
	}
}

func (vp *ValidatorPool) AddValidator(v *Validator) error {
	if vp.Validators == nil {
		return fmt.Errorf("validators is nil")
	}
	vp.Validators[v.Wallet.Address.Hex()] = v
	return nil
}

func (vp *ValidatorPool) SelectValidator(seed []byte) (common.Address, error) {
	if len(vp.Validators) == 0 {
		return common.Address{}, ErrNoValidators
	}

	var validatorList []*Validator
	for _, v := range vp.Validators {
		validatorList = append(validatorList, v)
	}

	hash := common.SHA256(seed)
	randNum := new(big.Int).SetBytes(hash.Bytes())

	selectedIndex := randNum.Mod(randNum, big.NewInt(int64(len(validatorList)))).Int64()
	selected := validatorList[selectedIndex]

	vp.SelectedValidator = selected

	return selected.Wallet.Address, nil
}
