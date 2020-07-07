package bob

import (
	"encoding/json"
	"fmt"
)

// E has address and value information
type E struct {
	A string `json:"a,omitempty" bson:"a,omitempty"`
	V uint32 `json:"v,omitempty" bson:"v,omitempty"`
	I uint8  `json:"i" bson:"i"`
	H string `json:"h,omitempty" bson:"h,omitempty"`
}

// Cell is a single OP_RETURN protocol
type Cell struct {
	H   string `json:"h,omitempty" bson:"h,omitempty"`
	B   string `json:"b,omitempty" bson:"b,omitempty"`
	S   string `json:"s,omitempty" bson:"s,omitempty"`
	I   uint8  `json:"i" bson:"i"`
	II  uint8  `json:"ii" bson:"ii"`
	Op  uint16 `json:"op,omitempty" bson:"op,omitempty"`
	Ops string `json:"ops,omitempty" bson:"ops,omitempty"`
}

// Input is a transaction input
type Input struct {
	I    uint8  `json:"i" bson:"i"`
	Tape []Tape `json:"tape" bson:"tape"`
	E    E      `json:"e" bson:"e"`
	Seq  uint32 `json:"seq" bson:"seq"`
}

// Tape is a tape
type Tape struct {
	Cell []Cell `json:"cell"`
	I    uint8  `json:"i"`
}

// Output is a transaction output
type Output struct {
	I    uint8  `json:"i"`
	Tape []Tape `json:"tape"`
	E    E      `json:"e,omitempty"`
}

// Blk containst the block info
type Blk struct {
	I uint32 `json:"i"`
}

// TxInfo conaints the transaction info
type TxInfo struct {
	H string `json:"h"`
}

// Tx is a BOB formatted Bitcoin transaction
type Tx struct {
	ID  string   `json:"_id"`
	Blk Blk      `json:"blk"`
	Tx  TxInfo   `json:"tx"`
	In  []Input  `json:"in"`
	Out []Output `json:"out"`
}

// New creates a new bob tx
func New() *Tx {
	return &Tx{}
}

// FromString takes a BOB formatted string
func (t *Tx) FromString(line string) {
	t.FromBytes([]byte(line))
}

// FromBytes takes a BOB formatted tx string as bytes
func (t *Tx) FromBytes(line []byte) error {
	if err := json.Unmarshal(line, &t); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}
