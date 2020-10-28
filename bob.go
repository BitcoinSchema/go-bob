package bob

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/libsv/libsv/script"
	"github.com/libsv/libsv/script/address"
	"github.com/libsv/libsv/transaction"
	"github.com/libsv/libsv/transaction/input"
	"github.com/libsv/libsv/transaction/output"
)

// OP_SWAP is 0x7c which is what "|" will be detected as
const asmProtocolDelimiter = "OP_SAWP"

// BobTxUnmarshal is a BOB formatted Bitcoin transaction that includes interfaces where types may change
type BobTxUnmarshal struct {
	ID   string            `json:"_id"`
	Blk  Blk               `json:"blk"`
	Tx   TxInfo            `json:"tx"`
	In   []Input           `json:"in"`
	Out  []OutputUnmarshal `json:"out"`
	Lock uint32            `json:"lock"`
}

// OutputUnmarshal is a transaction output
type OutputUnmarshal struct {
	I    uint8      `json:"i"`
	Tape []Tape     `json:"tape"`
	E    EUnmarshal `json:"e,omitempty"`
}

// EUnmarshal has address and value information
type EUnmarshal struct {
	A interface{} `json:"a,omitempty" bson:"a,omitempty"`
	V uint32      `json:"v,omitempty" bson:"v,omitempty"`
	I uint8       `json:"i" bson:"i"`
	H string      `json:"h,omitempty" bson:"h,omitempty"`
}

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
	LB  string `json:"lb,omitempty" bson:"lb,omitempty`
	S   string `json:"s,omitempty" bson:"s,omitempty"`
	LS  string `json:"ls,omitempty" bson:"ls,omitempty"`
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

// BobTx is a BOB formatted Bitcoin transaction
type BobTx struct {
	ID   string   `json:"_id"`
	Blk  Blk      `json:"blk"`
	Tx   TxInfo   `json:"tx"`
	In   []Input  `json:"in"`
	Out  []Output `json:"out"`
	Lock uint32   `json:"lock"`
}

// NewFromBytes creates a new BobTx from a NDJSON line representing a bob transaction, suck as returned by the bitbus 2 API
func NewFromBytes(line []byte) (bobTx *BobTx, err error) {
	bobTx = new(BobTx)
	err = bobTx.FromBytes(line)
	if err != nil {
		return nil, err
	}
	return bobTx, nil
}

// NewFromRawTxString creates a new BobTx from a hex encoded raw tx string
func NewFromRawTxString(rawTxString string) (bobTx *BobTx, err error) {
	bobTx = new(BobTx)
	err = bobTx.FromRawTxString(rawTxString)
	return
}

// NewFromString creates a new BobTx from a BOB formatted string
func NewFromString(line string) (bobTx *BobTx, err error) {
	bobTx = new(BobTx)
	err = bobTx.FromString(line)
	return
}

// NewFromTx creates a new BobTx from a libsv Transaction
func NewFromTx(tx *transaction.Transaction) (bobTx *BobTx, err error) {
	bobTx = new(BobTx)
	err = bobTx.FromTx(tx)
	return
}

// FromBytes takes a BOB formatted tx string as bytes
func (t *BobTx) FromBytes(line []byte) error {
	tu := new(BobTxUnmarshal)
	if err := json.Unmarshal(line, &tu); err != nil {
		return fmt.Errorf("Error parsing line: %s, %s", line, err)
	}

	// The out.E.A field can be either a boolean or a string
	// So we need to unmarshal into an interface, and fix the normal struct the user of this lib will work with (so they dont have to format the interface themselves)
	var fixedOuts []Output
	for _, out := range tu.Out {
		fixedOuts = append(fixedOuts, Output{
			I:    out.I,
			Tape: out.Tape,
			E: E{
				A: fmt.Sprintf("%s", out.E.A),
				V: out.E.V,
				I: out.E.I,
				H: out.E.H,
			},
		})
	}
	t.In = tu.In
	t.Blk = tu.Blk
	t.ID = tu.ID
	t.Out = fixedOuts
	t.Lock = tu.Lock
	t.Tx = tu.Tx

	// Check for missing hex values and supply them
	for outIdx, out := range t.Out {
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if len(cell.H) == 0 && len(cell.B) > 0 {
					// base 64 decode cell.B and encode it to hex string
					cellBytes, _ := base64.StdEncoding.DecodeString(cell.B)
					t.Out[outIdx].Tape[tapeIdx].Cell[cellIdx].H = hex.EncodeToString(cellBytes)
				}
			}
		}
	}
	for inIdx, in := range t.In {
		for tapeIdx, tape := range in.Tape {
			for cellIdx, cell := range tape.Cell {
				if len(cell.H) == 0 && len(cell.B) > 0 {
					// base 64 decode cell.B and encode it to hex string
					cellBytes, _ := base64.StdEncoding.DecodeString(cell.B)
					t.In[inIdx].Tape[tapeIdx].Cell[cellIdx].H = hex.EncodeToString(cellBytes)
				}
			}
		}
	}

	return nil
}

// FromRawTxString takes a hex encoded tx string
func (t *BobTx) FromRawTxString(rawTxString string) error {
	tx, err := transaction.NewFromString(rawTxString)
	if err != nil {
		return err
	}
	fmt.Printf("Calling with %+v", t)
	return t.FromTx(tx)
}

// FromString takes a BOB formatted string
func (t *BobTx) FromString(line string) (err error) {
	err = t.FromBytes([]byte(line))
	return
}

// FromTx takes a libsv.Transaction
func (t *BobTx) FromTx(tx *transaction.Transaction) error {

	// Set the transaction ID
	t.Tx.H = tx.GetTxID()

	// Set the inputs
	for inIdx, i := range tx.Inputs {

		bobInput := Input{
			I: uint8(inIdx),
			Tape: []Tape{{
				Cell: []Cell{{
					H: hex.EncodeToString(i.ToBytes(false)),
					B: base64.RawStdEncoding.EncodeToString(i.ToBytes(false)),
					S: i.String(),
				}},
				I: 0,
			}},
			E: E{
				H: i.PreviousTxID,
			},
		}

		t.In = append(t.In, bobInput)
	}

	// Process outputs
	for idxOut, o := range tx.Outputs {
		var adr string

		// Try to get a pubkeyhash (ignore fail when this is not a locking script)
		outPubKeyHash, _ := o.LockingScript.GetPublicKeyHash()
		if len(outPubKeyHash) > 0 {
			outAddress, err := address.NewFromPublicKeyHash(outPubKeyHash, true)
			if err != nil {
				return fmt.Errorf("Failed to get address from pubkeyhash %x: %s", outPubKeyHash, err)
			}
			adr = outAddress.AddressString
		}

		// Initialize out tapes and locking script asm
		asm, _ := o.LockingScript.ToASM()
		pushdatas := strings.Split(asm, " ")

		var outTapes []Tape
		bobOutput := Output{
			I:    uint8(idxOut),
			Tape: outTapes,
			E: E{
				A: adr,
			},
		}

		var currentTape Tape
		if len(pushdatas) > 0 {

			for pdIdx, pushdata := range pushdatas {
				pushdataBytes, _ := hex.DecodeString(pushdata)
				b64String := base64.StdEncoding.EncodeToString([]byte(pushdataBytes))

				if pushdata != asmProtocolDelimiter {
					currentTape.Cell = append(currentTape.Cell, Cell{
						B:  b64String,
						H:  pushdata,
						S:  string(pushdataBytes),
						I:  uint8(idxOut),
						II: uint8(pdIdx),
					})
				}
				// Note: OP_SWAP is 0x7c which is also ascii "|" which is our protocol separator. This is not used as OP_SWAP at all since this is in the script after the OP_FALSE
				if "OP_RETURN" == pushdata || asmProtocolDelimiter == pushdata {
					outTapes = append(outTapes, currentTape)
					currentTape = Tape{}
				}
			}
		}

		// Add the trailing tape
		outTapes = append(outTapes, currentTape)
		bobOutput.Tape = outTapes

		t.Out = append(t.Out, bobOutput)
	}

	return nil
}

// ToRawTxString converts the BOBTx to a libsv.transaction, and outputs the raw hex
func (t *BobTx) ToRawTxString() (string, error) {
	tx, err := t.ToTx()
	if err != nil {
		return "", err
	}
	return tx.ToString(), nil
}

// ToString returns a json string of bobTx
func (t *BobTx) ToString() (string, error) {
	// Create JSON from the instance data.
	// ... Ignore errors.
	b, err := json.Marshal(t)
	// Convert bytes to string.
	return string(b), err

}

// ToTx returns a transaction.Transaction
func (t *BobTx) ToTx() (*transaction.Transaction, error) {
	tx := transaction.New()

	tx.Locktime = t.Lock

	for _, in := range t.In {

		if len(in.Tape) == 0 || len(in.Tape[0].Cell) == 0 {
			return nil, fmt.Errorf("Failed to process inputs. More tapes or cells than expected. %+v", in.Tape)
		}

		prevTxScript, _ := script.NewP2PKHFromAddress(fmt.Sprintf("%s", in.E.A))

		var scriptAsm []string
		for _, cell := range in.Tape[0].Cell {
			cellData := cell.H
			scriptAsm = append(scriptAsm, cellData)
		}

		builtUnlockScript, err := script.NewFromASM(strings.Join(scriptAsm, " "))
		if err != nil {
			log.Println("Failed to get script from asm", scriptAsm, err)
		}

		// add inputs
		i := &input.Input{
			PreviousTxID:       in.E.H,
			PreviousTxOutIndex: uint32(in.E.I),
			PreviousTxSatoshis: uint64(in.E.V),
			PreviousTxScript:   prevTxScript,
			UnlockingScript:    builtUnlockScript,
			SequenceNumber:     in.Seq,
		}

		tx.AddInput(i)
	}

	// add outputs
	for _, out := range t.Out {
		// Build the locking script
		var lockScriptAsm []string
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if cellIdx == 0 && tapeIdx > 1 {
					// add the separator back in
					lockScriptAsm = append(lockScriptAsm, asmProtocolDelimiter)
				}

				if len(cell.H) > 0 {
					lockScriptAsm = append(lockScriptAsm, cell.H)
				} else if len(cell.Ops) > 0 {
					lockScriptAsm = append(lockScriptAsm, cell.Ops)
				}
			}
		}

		lockingScript, _ := script.NewFromASM(strings.Join(lockScriptAsm, " "))
		o := &output.Output{
			Satoshis:      uint64(out.E.V),
			LockingScript: lockingScript,
		}

		tx.AddOutput(o)
	}

	return tx, nil
}
