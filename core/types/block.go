package types

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/PulseCoinOrg/nexacoin/common"
)

var (
	NoTxHash = common.SHA256([]byte("000000000000000000000000000000"))
)

type Block struct {
	Time         int64
	Hash         common.Hash
	ParentHash   common.Hash
	UncleHash    common.Hash // A hash of a block that is valid but not chosen to be apart of the chain
	Transactions []*Transaction
	TxHash       common.Hash
	// TODO Gas uint64 add this
}

func NewBlock(time int64, parentHash common.Hash, transactions []*Transaction) *Block {
	block := &Block{
		Time:         time,
		ParentHash:   parentHash,
		Transactions: transactions,
	}
	block.Hash = common.SHA256(block.BytesStream())
	if len(transactions) == 0 {
		block.TxHash = NoTxHash
	}
	return block
}

func (b *Block) BytesStream() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

func DecodeBlockBytesStream(data []byte) *Block {
	var block Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&block); err != nil {
		fmt.Println(err)
	}
	return &block
}
