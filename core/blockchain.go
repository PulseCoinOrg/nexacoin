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

package core

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/PulseCoinOrg/nexacoin/common"
	"github.com/PulseCoinOrg/nexacoin/core/types"
	"github.com/PulseCoinOrg/nexacoin/nexadb/leveldb"
)

var (
	// TODO add support for other operating systems & change the path to ~/home/...
	ChainDiskPath = "./chaindb-output"
)

// removes the chain disk leveldb folder from systsem
func DeleteDiskFolder() error {
	err := os.Remove(ChainDiskPath)
	if err != nil {
		return err
	}
	return nil
}

type BlockChain struct {
	Database     *leveldb.Database
	Height       uint64
	Sane         bool
	LastBlock    *types.Block
	BlocksMemory map[common.Hash]*types.Block
	Validators   *ValidatorPool
}

func NewChain() (*BlockChain, error) {
	db, err := leveldb.New(ChainDiskPath)
	if db == nil {
		return nil, ErrChainDatabaseClosed
	}
	if err != nil {
		return nil, err
	}
	return &BlockChain{
		Database:     db,
		BlocksMemory: make(map[common.Hash]*types.Block),
		Validators:   NewValidatorPool(),
	}, nil
}

// returns the first element in the chain from leveldb
func (chain *BlockChain) First() (*types.Block, error) {
	_, value, err := chain.Database.First()
	if err != nil {
		return nil, err
	}
	block := types.DecodeBlockBytesStream(value)
	return block, nil
}

// returns the last element in the chain from leveldb
func (chain *BlockChain) Last() (*types.Block, error) {
	_, value, err := chain.Database.Last()
	if err != nil {
		return nil, err
	}
	block := types.DecodeBlockBytesStream(value)
	chain.LastBlock = block
	return block, nil
}

func (chain *BlockChain) Previous() (*types.Block, error) {
	_, value, err := chain.Database.Previous()
	if err != nil {
		return nil, err
	}
	block := types.DecodeBlockBytesStream(value)
	chain.LastBlock = block
	return block, nil
}

// returns the second to last element in the chain from leveldb
func (chain *BlockChain) Insert(b *types.Block) error {
	if err := chain.Database.Put(b.Hash.Bytes(), b.BytesStream()); err != nil {
		return ErrBlockChainInsertFailed
	}
	if chain.BlocksMemory == nil {
		return ErrBlockChainInsertFailed
	}
	chain.BlocksMemory[b.Hash] = b
	return nil
}

// this is equivilent to a BlockByHash function
// retrieves a block by a given hash string
func (chain *BlockChain) LocateBlock(hash string) *types.Block {
	for _, block := range chain.BlocksMemory {
		if block.Hash.Hex() == hash {
			return block
		}
	}
	return nil
}

// checks if the chain is sane (AKA valid)
func (chain *BlockChain) SanityCheck() bool {
	lastBlock, err := chain.Last()
	if err != nil || lastBlock == nil {
		chain.Sane = false
		return false
	}

	current := lastBlock
	for current.ParentHash != (common.Hash{}) {
		parent, ok := chain.BlocksMemory[current.ParentHash]
		if !ok {
			chain.Sane = false
			return false
		}

		if current.ParentHash.Hex() != parent.Hash.Hex() {
			chain.Sane = false
			return false
		}

		current = parent
	}

	chain.Sane = true
	return true
}

// validates the blocks from an electoral system.
// a validator is picked from a pool of validators for now
// a validator must stake a minimum of [x]nex to participate
// if someone tries tricking the network, all of the staked crypto will be lost
// TODO 'x' nex must be a reasonable amount as we start, but not reasonable enough that
// everyone can participate.
func (chain *BlockChain) pickValidator() (*Validator, error) {
	prevBlock, err := chain.Previous()
	if err != nil {
		return nil, fmt.Errorf("error fetching second to last block")
	}

	chain.LastBlock = prevBlock

	if chain.LastBlock == nil {
		return nil, fmt.Errorf("LastBlock is nil")
	}

	address, err := chain.Validators.SelectValidator(chain.LastBlock.Hash.Bytes())
	if err != nil {
		return nil, ErrBlockChainValidatorSelectFailed
	}

	validator, ok := chain.Validators.Validators[address.Hex()]
	if !ok || validator == nil {
		return nil, fmt.Errorf("validator with address %s not found", address.Hex())
	}

	chain.Validators.SelectedValidator = validator
	return validator, nil
}

// picks a validator and uses them to validate the latest block in the chain
func (chain *BlockChain) ValidateLastBlock() bool {
	validator, err := chain.pickValidator()
	if err != nil {
		slog.Error("Failed to pick validator", "err", err)
		return false
	}

	lastBlock, err := chain.Last()
	if err != nil || lastBlock == nil {
		slog.Error("Failed to load last block", "err", err)
		return false
	}

	addr, err := validator.GetValidatorAddress()
	if addr == "" {
		slog.Error("Failed to get validator address", "err", err)
	}
	slog.Info("validator has been chosen", "addr", addr)

	return validator.ValidateBlock(lastBlock)
}
