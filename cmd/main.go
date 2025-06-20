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

package main

import (
	"log/slog"
	"time"

	"github.com/PulseCoinOrg/nexacoin/core"
	"github.com/PulseCoinOrg/nexacoin/core/types"
	"github.com/PulseCoinOrg/nexacoin/wallet"
)

func main() {
	w, err := wallet.New()
	Handle(err)

	err = w.SaveDisk()
	Handle(err)

	chain, err := core.NewChain()
	Handle(err)

	v, err := core.NewValidator()
	Handle(err)

	err = chain.Validators.AddValidator(v)
	Handle(err)

	block1 := types.NewBlock(time.Now().Unix(), core.GenesisParentHash, []*types.Transaction{})
	err = chain.Insert(block1)
	Handle(err)

	block2 := types.NewBlock(time.Now().Unix(), block1.Hash, []*types.Transaction{})
	err = chain.Insert(block2)
	Handle(err)

	block3 := types.NewBlock(time.Now().Unix(), block2.Hash, []*types.Transaction{})
	err = chain.Insert(block3)
	Handle(err)

	valid := chain.ValidateLastBlock()
	if !valid {
		slog.Error("chain validator has found an invalid block")
	} else {
		slog.Info("chain validated using validator")
	}
}

func Handle(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}
