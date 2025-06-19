package types

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/PulseCoinOrg/nexacoin/common"
)

type Transaction struct {
	Time      int64
	Fee       uint64
	Sender    common.Address
	Recipient common.Address
	Amount    int64
	Hash      common.Hash
}

func NewTx(
	time int64,
	sender common.Address,
	recipient common.Address,
	amount int64,
) *Transaction {
	return &Transaction{
		Time:      time,
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}

// converts the transaction into bytes
func (tx *Transaction) BytesStream() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

// converts transaction bytes into a transaction
func DecodeTxBytesStream(data []byte) *Transaction {
	var tx Transaction
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&tx); err != nil {
		fmt.Println(err)
	}
	return &tx
}
