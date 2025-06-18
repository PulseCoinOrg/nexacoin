package core

import (
	"os"

	"github.com/PulseCoinOrg/nexacoin/common"
	"github.com/PulseCoinOrg/nexacoin/core/types"
	"github.com/PulseCoinOrg/nexacoin/nexadb/leveldb"
)

var (
	// TODO add support for other operating systems & change the path to ~/home/...
	ChainDiskFile = "./chaindb-output"
)

func DeleteDiskFile() error {
	err := os.Remove(ChainDiskFile)
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
}

func NewChain() (*BlockChain, error) {
	db, err := leveldb.New(ChainDiskFile)
	if db == nil {
		return nil, ErrChainDatabaseClosed
	}
	if err != nil {
		return nil, err
	}
	return &BlockChain{
		Database:     db,
		BlocksMemory: make(map[common.Hash]*types.Block),
	}, nil
}

func (chain *BlockChain) First() (*types.Block, error) {
	_, value, err := chain.Database.First()
	if err != nil {
		return nil, err
	}
	block := types.DecodeBlockBytesStream(value)
	return block, nil
}

func (chain *BlockChain) Last() (*types.Block, error) {
	_, value, err := chain.Database.Last()
	if err != nil {
		return nil, err
	}
	block := types.DecodeBlockBytesStream(value)
	return block, nil
}

func (chain *BlockChain) Previous() (*types.Block, error) {
	_, value, err := chain.Database.Previous()
	if err != nil {
		return nil, err
	}
	block := types.DecodeBlockBytesStream(value)
	return block, nil
}

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

func (chain *BlockChain) SanityCheck() bool {
	lastBlock, err := chain.Last()
	if err != nil || lastBlock == nil {
		chain.Sane = false
		return false
	}

	// Check if block exists in memory
	current := lastBlock
	for current.ParentHash != (common.Hash{}) { // walk until genesis
		parent, ok := chain.BlocksMemory[current.ParentHash]
		if !ok {
			chain.Sane = false
			return false
		}

		// Validate linkage
		if current.ParentHash.Hex() != parent.Hash.Hex() {
			chain.Sane = false
			return false
		}

		current = parent
	}

	chain.Sane = true
	return true
}
