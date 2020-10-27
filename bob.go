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

// Tx is a BOB formatted Bitcoin transaction
type BobTx struct {
	ID  string   `json:"_id"`
	Blk Blk      `json:"blk"`
	Tx  TxInfo   `json:"tx"`
	In  []Input  `json:"in"`
	Out []Output `json:"out"`
}

// New creates a new bob tx
func New() *BobTx {
	return &BobTx{}
}

// FromString takes a BOB formatted string
func (t *BobTx) FromString(line string) error {
	err := t.FromBytes([]byte(line))
	if err != nil {
		return err
	}
	return nil
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

		outPubKeyHash, err := o.LockingScript.GetPublicKeyHash()
		if err != nil {
			log.Printf("oh no 1 %x: %s", outPubKeyHash, err)
		}
		if len(outPubKeyHash) > 0 {
			outAddress, err := address.NewFromPublicKeyHash(outPubKeyHash, true)
			if err != nil {
				return fmt.Errorf("oh no 2 %x: %s", outPubKeyHash, err)
			}
			adr = outAddress.AddressString
		}

		// Inspect OP_RETURN data
		// Find OP_FALSE + OP_RETURN -OR- OP_RETURN
		// Those go in the first tape
		asm, _ := o.LockingScript.ToASM()
		fmt.Println("Locking script", asm)
		pushdatas := strings.Split(asm, " ")

		var outTapes []Tape
		// TODO: Find | pushdata, split into tapes
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
				log.Printf("Pushdata %d: %s \n", pdIdx, pushdata)
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
					log.Println("End of tape detected", pushdata)
					outTapes = append(outTapes, currentTape)
					currentTape = Tape{}
				}
			}

			//

		}
		// Add the trailing tape
		outTapes = append(outTapes, currentTape)

		bobOutput.Tape = outTapes

		t.Out = append(t.Out, bobOutput)
	}

	return nil
}

// ToTx returns a transaction.Transaction
func (t *BobTx) ToTx() (*transaction.Transaction, error) {

	tx := transaction.New()

	log.Printf("Processing bob tx. Ins: %d, Outs: %d ~~~~~ %+v", len(t.In), len(t.Out), t)

	for _, in := range t.In {
		if len(in.Tape) == 0 || len(in.Tape[0].Cell) == 0 {
			return nil, fmt.Errorf("Failed to process inputs. More tapes or cells than expected. %+v", in.Tape)
		}
		log.Println("In Tapes", len(in.Tape))
		builtPrevTxScript := in.Tape[0].Cell[0].H
		prevTxScript, _ := script.NewFromHexString(builtPrevTxScript)

		// add inputs
		i := &input.Input{
			PreviousTxID:       in.E.H,
			PreviousTxOutIndex: uint32(in.E.I),
			PreviousTxSatoshis: uint64(in.E.V),
			PreviousTxScript:   prevTxScript,
		}
		tx.AddInput(i)
	}

	// add outputs
	for _, out := range t.Out {
		// Build the locking script
		var lockScriptAsm []string //out.Tape[0].Cell[0].B
		log.Println("Out Tapes", len(out.Tape))
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if cellIdx == 0 && tapeIdx > 1 {
					// add the separator back in
					log.Println("Appending", asmProtocolDelimiter)
					lockScriptAsm = append(lockScriptAsm, asmProtocolDelimiter)
				}
				log.Println("Appending", cell.H)

				if len(cell.H) > 0 {
					lockScriptAsm = append(lockScriptAsm, cell.H)
				} else if len(cell.Ops) > 0 {
					lockScriptAsm = append(lockScriptAsm, cell.Ops)
				}
			}
		}

		log.Printf("Assembed asm %s", lockScriptAsm)
		lockingScript, _ := script.NewFromASM(strings.Join(lockScriptAsm, " "))

		o := &output.Output{
			Satoshis:      uint64(out.E.V),
			LockingScript: lockingScript,
		}

		tx.AddOutput(o)
	}

	return tx, nil
}

// ToRawTxString converts the BOBTx to a libsv.transaction, and outputs the raw hex
func (t *BobTx) ToRawTxString() (string, error) {
	tx, err := t.ToTx()
	if err != nil {
		return "", err
	}
	return tx.ToString(), nil
}

// FromBytes takes a BOB formatted tx string as bytes
func (t *BobTx) FromBytes(line []byte) error {
	if err := json.Unmarshal(line, &t); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}
