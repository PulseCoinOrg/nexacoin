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
	Height       int
	// Nonce        int
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

// converts the block into bytes
func (b *Block) BytesStream() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

// converts the bytes of a block into a block
func DecodeBlockBytesStream(data []byte) *Block {
	var block Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&block); err != nil {
		fmt.Println(err)
	}
	return &block
}
