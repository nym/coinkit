package data

import (
	"fmt"

	"github.com/lacker/coinkit/util"
)

// UpdateOperation is used to alter the contents of a document that is already stored.
type UpdateOperation struct {
	// Who is updating the document. Must be the owner
	Signer string

	// The sequence number for this operation
	Sequence uint32

	// The data to update the document with.
	Data *JSONObject

	// How much the updater is willing to pay to send this operation through
	Fee uint64
}

func (op *UpdateOperation) String() string {
	return fmt.Sprintf("update owner=%s, data=%s", util.Shorten(op.Signer), op.Data)
}

func (op *UpdateOperation) OperationType() string {
	return "Update"
}

func (op *UpdateOperation) GetSigner() string {
	return op.Signer
}

func (op *UpdateOperation) GetFee() uint64 {
	return op.Fee
}

func (op *UpdateOperation) GetSequence() uint32 {
	return op.Sequence
}

// TODO: should this do something?
func (op *UpdateOperation) Verify() bool {
	return true
}

// Works with MakeTestCreateOperation to change the value
func MakeTestUpdateOperation(n int, sequence int) *SignedOperation {
	mint := util.NewKeyPairFromSecretPhrase("mint")
	data := NewEmptyJSONObject()
	data.Set("foo", sequence)
	op := &UpdateOperation{
		Signer:   mint.PublicKey().String(),
		Sequence: uint32(sequence),
		Data:     data,
		Fee:      0,
	}
	return NewSignedOperation(op, mint)
}

func init() {
	RegisterOperationType(&UpdateOperation{})
}
