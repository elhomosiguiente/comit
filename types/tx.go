package types

import (
	"encoding/json"
	"fmt"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	CreateAccountTx = 0x01
	RemoveAccountTx = 0x02
	SubmitTx        = 0x03
	ResolveTx       = 0x04

	// Errors
	ErrUnexpectedData       = 1311
	ErrAccountAlreadyExists = 11311
)

type TxInput struct {
	Address   []byte
	Sequence  int
	Signature crypto.Signature
	PubKey    crypto.PubKey
}

func (txIn TxInput) ValidateBasic() tmsp.Result {
	if len(txIn.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if txIn.Sequence <= 0 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Sequence must be greater than 0")
	}
	if txIn.Sequence == 1 && txIn.PubKey == nil {
		return tmsp.ErrBaseInvalidInput.AppendLog("PubKey must be present when Sequence == 1")
	}
	if txIn.Sequence > 1 && txIn.PubKey != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog("PubKey must be nil when Sequence > 1")
	}
	return tmsp.OK
}

func (txIn TxInput) String() string {
	return fmt.Sprintf("TxInput{%X,%v,%v,%v}", txIn.Address, txIn.Sequence, txIn.Signature, txIn.PubKey)
}

type Tx struct {
	Type  byte
	Input TxInput
	Data  []byte
}

func (tx *Tx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Input.Signature
	tx.Input.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Input.Signature = sig
	return signBytes
}

func (tx *Tx) SetAccount(addr []byte) {
	var pubKey crypto.PubKeyEd25519
	tx.Input.Address = addr
	copy(pubKey[:], addr[:])
	tx.Input.PubKey = pubKey
}

func (tx *Tx) SetSignature(sigBytes []byte) {
	var sig crypto.SignatureEd25519
	copy(sig[:], sigBytes[:])
	tx.Input.Signature = sig
}

func (tx *Tx) String() string {
	return fmt.Sprintf("Tx{%v %v %X}", tx.Type, tx.Input, tx.Data)
}

func TxID(chainID string, tx Tx) []byte {
	signBytes := tx.SignBytes(chainID)
	return wire.BinaryRipemd160(signBytes)
}

func jsonEscape(str string) string {
	escapedBytes, err := json.Marshal(str)
	if err != nil {
		PanicSanity(Fmt("Error json-escaping a string", str))
	}
	return string(escapedBytes)
}