// Package bob is a library for working with BOB formatted transactions
//
// Specs: https://bob.planaria.network/
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package bob

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitcoinschema/go-bpu"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/template/p2pkh"
)

// Protocol delimiter constants
// OP_SWAP = 0x7c = 124 = "|"
const (
	ProtocolDelimiterAsm  = "OP_SWAP"
	ProtocolDelimiterInt  = 0x7c
	ProtocolDelimiterByte = byte(ProtocolDelimiterInt)
	ProtocolDelimiter     = string(rune(ProtocolDelimiterInt))
)

// Tx is a BOB formatted Bitcoin transaction
//
// DO NOT CHANGE ORDER - aligned for memory optimization (malign)
type Tx struct {
	bpu.Tx
}

// used by bpu.Parse to determine if the parsing should be shallow or deep
// always using shallow since it covers 99.99% of cases and eliminates
// bottlenecking on txs with lots of pushdatas (like complex sCrypt contracts)
var shallowMode = bpu.Shallow

// NewFromBytes creates a new BOB Tx from a NDJSON line representing a BOB transaction,
// as returned by the bitbus 2 API
func NewFromBytes(line []byte) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromBytes(line)
	return
}

// NewFromRawTxString creates a new BobTx from a hex encoded raw tx string
func NewFromRawTxString(rawTxString string) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromRawTxString(rawTxString)
	return
}

// NewFromString creates a new BobTx from a BOB formatted string
func NewFromString(line string) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromString(line)
	return
}

// NewFromTx creates a new BobTx from a libsv Transaction
func NewFromTx(tx *transaction.Transaction) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromTx(tx)
	if err != nil {
		return nil, err
	}
	return
}

// FromBytes takes a BOB formatted tx string as bytes
func (t *Tx) FromBytes(line []byte) error {
	tu := new(bpu.Tx)
	if err := json.Unmarshal(line, &tu); err != nil {
		return fmt.Errorf("error parsing line: %v, %w", line, err)
	}

	// The out.E.A field can be either an address or "false"
	fixedOuts := make([]bpu.Output, 0)
	for _, out := range tu.Out {
		var address string
		if out.E.A != nil {
			address = *out.E.A
		}
		fixedOuts = append(fixedOuts, bpu.Output{
			XPut: bpu.XPut{
				I:    out.I,
				Tape: out.Tape,
				E: bpu.E{
					A: &address,
					V: out.E.V,
					I: out.E.I,
					H: out.E.H,
				},
			},
		})
	}
	t.Blk = tu.Blk
	t.ID = tu.ID
	t.In = tu.In
	t.Lock = tu.Lock
	t.Out = fixedOuts

	t.Tx.Tx = tu.Tx

	// Check for missing hex values and supply them
	for outIdx, out := range t.Out {
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if cell.H == nil && cell.B != nil && len(*cell.B) > 0 {
					// base 64 decode cell.B and encode it to hex string
					cellBytes, err := base64.StdEncoding.DecodeString(*cell.B)
					if err != nil {
						return err
					}
					var hexStr = hex.EncodeToString(cellBytes)
					t.Out[outIdx].Tape[tapeIdx].Cell[cellIdx].H = &hexStr
				}
			}
		}
	}
	for inIdx, in := range t.In {
		for tapeIdx, tape := range in.Tape {
			for cellIdx, cell := range tape.Cell {
				if cell.H == nil && cell.B != nil && len(*cell.B) > 0 {
					// base 64 decode cell.B and encode it to hex string
					cellBytes, err := base64.StdEncoding.DecodeString(*cell.B)
					if err != nil {
						return err
					}
					hexStr := hex.EncodeToString(cellBytes)
					t.In[inIdx].Tape[tapeIdx].Cell[cellIdx].H = &hexStr
				}
			}
		}
	}

	return nil
}

// FromRawTxString takes a hex encoded tx string
func (t *Tx) FromRawTxString(rawTxString string) (err error) {

	var separator = "|"
	var l = bpu.IncludeL
	var opReturn = uint8(106)
	var opFalse = uint8(0)
	var splitConfig = []bpu.SplitConfig{
		{
			Token: &bpu.Token{
				Op: &opReturn,
			},
			Include: &l,
		},
		{
			Token: &bpu.Token{
				Op: &opFalse,
			},
			Include: &l,
			Require: &opReturn,
		},
		{
			Token: &bpu.Token{
				S: &separator,
			},
			Require: &opReturn,
		},
	}

	bpuTx, err := bpu.Parse(bpu.ParseConfig{RawTxHex: &rawTxString, SplitConfig: splitConfig, Mode: &shallowMode})
	if bpuTx != nil {
		t.Tx = *bpuTx
	}

	return
}

// FromString takes a BOB formatted string
func (t *Tx) FromString(line string) (err error) {
	err = t.FromBytes([]byte(line))
	return
}

// FromTx takes a bt.Tx
func (t *Tx) FromTx(tx *transaction.Transaction) error {

	if tx == nil {
		return fmt.Errorf("Tx must be set")
	}
	var separator = "|"
	var l = bpu.IncludeL
	var opReturn = uint8(106)
	var splitConfig = []bpu.SplitConfig{
		{
			Token: &bpu.Token{
				Op: &opReturn,
			},
			Include: &l,
		},
		{
			Token: &bpu.Token{
				S: &separator,
			},
			Require: &opReturn,
		},
	}

	bpuTx, err := bpu.Parse(bpu.ParseConfig{Tx: tx, SplitConfig: splitConfig, Mode: &shallowMode})
	if err != nil {
		return err
	}
	if bpuTx != nil {
		t.Tx = *bpuTx
	}
	return nil
}

// ToRawTxString converts the BOBTx to a libsv.transaction, and outputs the raw hex
func (t *Tx) ToRawTxString() (string, error) {
	tx, err := t.ToTx()
	if err != nil {
		return "", err
	}
	return tx.String(), nil
}

// ToString returns a json string of bobTx
func (t *Tx) ToString() (string, error) {
	// Create JSON from the instance data.
	b, err := json.Marshal(t)
	return string(b), err

}

// ToTx returns a bt.Tx
func (t *Tx) ToTx() (*transaction.Transaction, error) {
	tx := transaction.NewTransaction()

	tx.LockTime = t.Lock

	for _, in := range t.In {

		if len(in.Tape) == 0 || len(in.Tape[0].Cell) == 0 {
			return nil, fmt.Errorf("failed to process inputs. More tapes or cells than expected. %+v", in.Tape)
		}

		add, err := script.NewAddressFromString(*in.E.A)
		if err != nil {
			return nil, err
		}
		prevTxScript, _ := p2pkh.Lock(add)

		var scriptAsm []string
		// TODO: This will break if there is ever a bpu splitter present in inputs
		for _, cell := range in.Tape[0].Cell {
			cellData := *cell.H
			scriptAsm = append(scriptAsm, cellData)
		}

		builtUnlockScript, err := script.NewFromASM(strings.Join(scriptAsm, " "))
		if err != nil {
			return nil, fmt.Errorf("failed to get script from asm: %v error: %w", scriptAsm, err)
		}
		v := uint64(0)
		if in.E.V != nil {
			v = *in.E.V
		}

		// add inputs
		i := &transaction.TransactionInput{
			SourceTxOutIndex: in.E.I, // TODO: This might be getting set incorrectly?
			UnlockingScript:  builtUnlockScript,
			SequenceNumber:   in.Seq,
		}
		i.SourceTXID, _ = chainhash.NewHashFromHex(*in.E.H)
		i.SetSourceTxOutput(&transaction.TransactionOutput{
			Satoshis:      v,
			LockingScript: prevTxScript,
		})

		tx.Inputs = append(tx.Inputs, i) // AddInput(i)
	}

	// add outputs
	for _, out := range t.Out {
		// Build the locking script
		var lockScriptAsm []string
		for tapeIdx, tape := range out.Tape {
			for cellIdx, cell := range tape.Cell {
				if cellIdx == 0 && tapeIdx > 1 {
					// add the separator back in
					lockScriptAsm = append(lockScriptAsm, ProtocolDelimiterAsm)
				}

				if cell.H != nil {
					lockScriptAsm = append(lockScriptAsm, *cell.H)
				} else if cell.Ops != nil {
					lockScriptAsm = append(lockScriptAsm, *cell.Ops)
				}
			}
		}

		lockingScript, _ := script.NewFromASM(strings.Join(lockScriptAsm, " "))
		o := &transaction.TransactionOutput{
			Satoshis:      *out.E.V,
			LockingScript: lockingScript,
		}

		tx.AddOutput(o)
	}

	return tx, nil
}
