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

	"github.com/bitcoinschema/go-bob/util"
	"github.com/bitcoinschema/go-bpu"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
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
	bpu.BpuTx
}

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
func NewFromTx(tx *bt.Tx) (bobTx *Tx, err error) {
	bobTx = new(Tx)
	err = bobTx.FromTx(tx)
	return
}

// FromBytes takes a BOB formatted tx string as bytes
func (t *Tx) FromBytes(line []byte) error {
	tu := new(bpu.BpuTx)
	if err := json.Unmarshal(line, &tu); err != nil {
		return fmt.Errorf("error parsing line: %v, %w", line, err)
	}

	// The out.E.A field can be either an address or "false"
	fixedOuts := make([]bpu.Output, 0)
	for _, out := range tu.Out {
		address := fmt.Sprintf("%s", *out.E.A)
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
	t.Tx = tu.Tx

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

	var seperator = "|"
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
		},
		{
			Token: &bpu.Token{
				S: &seperator,
			},
		},
	}

	bpuTx, err := bpu.Parse(bpu.ParseConfig{RawTxHex: rawTxString, SplitConfig: splitConfig})
	t.BpuTx = *bpuTx

	return
}

// FromString takes a BOB formatted string
func (t *Tx) FromString(line string) (err error) {
	err = t.FromBytes([]byte(line))
	return
}

// FromTx takes a bt.Tx
func (t *Tx) FromTx(tx *bt.Tx) error {

	// Set the transaction ID
	t.Tx.H = tx.TxID()

	// Set the inputs
	for inIdx, i := range tx.Inputs {

		cellHex := hex.EncodeToString(i.Bytes(false))
		cellB64 := base64.RawStdEncoding.EncodeToString(i.Bytes(false))
		cellStr := i.String()
		txid := hex.EncodeToString(i.PreviousTxID())

		bobInput := bpu.Input{
			XPut: bpu.XPut{
				I: uint8(inIdx),
				Tape: []bpu.Tape{{
					Cell: []bpu.Cell{{
						H: &cellHex,
						B: &cellB64,
						S: &cellStr,
					}},
					I: 0,
				}},
				E: bpu.E{
					H: &txid,
				},
			},
		}

		t.In = append(t.In, bobInput)
	}

	// Process outputs
	for idxOut, o := range tx.Outputs {
		var adr string

		// Try to get a pub_key hash (ignore fail when this is not a locking script)
		outPubKeyHash, _ := o.LockingScript.PublicKeyHash()
		if len(outPubKeyHash) > 0 {
			outAddress, err := bscript.NewAddressFromPublicKeyHash(outPubKeyHash, true)
			if err != nil {
				return fmt.Errorf("failed to get address from pubkeyhash %x: %w", outPubKeyHash, err)
			}
			adr = outAddress.AddressString
		}

		// Initialize out tapes and locking script asm
		asm, err := o.LockingScript.ToASM()
		if err != nil {
			return err
		}

		pushDatas := strings.Split(asm, " ")

		var outTapes []bpu.Tape
		bobOutput := bpu.Output{
			XPut: bpu.XPut{
				I:    uint8(idxOut),
				Tape: outTapes,
				E: bpu.E{
					A: &adr,
				},
			},
		}

		var opTape bpu.Tape
		var currentTape bpu.Tape
		var opOffset = 0
		if len(pushDatas) > 0 {

			// Check for OP_RETURN or OP_FALSE + OP_RETURN
			// Look for OP_FALSE OP_RETURN or just OP_RETURN and separate into a cell collection
			if len(pushDatas[0]) > 0 {
				if pushDatas[0] == "OP_FALSE" || pushDatas[0] == "0" {
					// OP_FALSE in position 0
					var op = uint8(bscript.OpFALSE)
					var ops = "OP_FALSE"
					opTape.Cell = append(opTape.Cell, bpu.Cell{
						Op:  &op,
						Ops: &ops,
						I:   uint8(idxOut),
						II:  uint8(0),
					})
					opOffset += 1
					// Check for OP_RETURN
					if len(pushDatas[1]) > 0 && pushDatas[1] == "OP_RETURN" {
						// OP_FALSE OP_RETURN
						var op = uint8(bscript.OpRETURN)
						var ops = "OP_RETURN"
						opTape.Cell = append(opTape.Cell, bpu.Cell{
							Op:  &op,
							Ops: &ops,
							I:   uint8(idxOut),
							II:  uint8(1),
						})
						// pull them out into their own cell collection
						outTapes = append(outTapes, opTape)
						opOffset += 1

					}
				} else if len(pushDatas[0]) > 0 && pushDatas[0] == "OP_RETURN" {
					var op = uint8(bscript.OpRETURN)
					var ops = "OP_RETURN"
					opTape.Cell = append(opTape.Cell, bpu.Cell{
						Op:  &op,
						Ops: &ops,
						I:   uint8(idxOut),
						II:  uint8(0),
					})
					opOffset += 1

					// OP_RETURN in position 0
				}
				if opOffset > 0 {
					outTapes = append(outTapes, opTape)
				}
			}
			for pdIdx, pushData := range pushDatas {
				if pdIdx < opOffset {
					continue
				}
				// Ignore error if it fails, use empty
				pushDataBytes, _ := hex.DecodeString(pushData)
				b64String := base64.StdEncoding.EncodeToString(pushDataBytes)
				var pushDataString = string(pushDataBytes)

				// assume the pushdata is a chunk of
				pushDataHex := pushData
				var op uint8
				var ops string
				// asm is being put into the hex field - need to convert back to hex from opcodes if they exist in here
				if pushDataByte, ok := util.OpCodeStrings[pushData]; ok {
					// this pushdata is a valid opcode
					pushDataHex = hex.EncodeToString([]byte{pushDataByte})
					op = uint8(pushDataByte)
					ops = pushData
				}
				if pushData != ProtocolDelimiterAsm {
					currentTape.Cell = append(currentTape.Cell, bpu.Cell{
						Op:  &op,
						Ops: &ops,
						B:   &b64String,
						H:   &pushDataHex,
						S:   &pushDataString,
						I:   uint8(idxOut),
						II:  uint8(pdIdx - opOffset),
					})
				}
				// Note: OP_SWAP is 0x7c which is also ascii "|" which is our protocol separator.
				// This is not used as OP_SWAP at all since this is in the script after the OP_FALSE
				// if "OP_RETURN" == pushData || ProtocolDelimiterAsm == pushData {
				if ProtocolDelimiterAsm == pushData {
					outTapes = append(outTapes, currentTape)
					currentTape = bpu.Tape{}
					opOffset = 0
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
func (t *Tx) ToTx() (*bt.Tx, error) {
	tx := bt.NewTx()

	tx.LockTime = t.Lock

	for _, in := range t.In {

		if len(in.Tape) == 0 || len(in.Tape[0].Cell) == 0 {
			return nil, fmt.Errorf("failed to process inputs. More tapes or cells than expected. %+v", in.Tape)
		}

		prevTxScript, _ := bscript.NewP2PKHFromAddress(*in.E.A)

		var scriptAsm []string
		// TODO: This will break if there is ever a bpu splitter present in inputs
		for _, cell := range in.Tape[0].Cell {
			cellData := *cell.H
			scriptAsm = append(scriptAsm, cellData)
		}

		builtUnlockScript, err := bscript.NewFromASM(strings.Join(scriptAsm, " "))
		if err != nil {
			return nil, fmt.Errorf("failed to get script from asm: %v error: %w", scriptAsm, err)
		}
		v := uint64(0)
		if in.E.V != nil {
			v = *in.E.V
		}

		// add inputs
		i := &bt.Input{
			PreviousTxOutIndex: in.E.I, // TODO: This might be getting set incorrectly?
			PreviousTxSatoshis: v,
			PreviousTxScript:   prevTxScript,
			UnlockingScript:    builtUnlockScript,
			SequenceNumber:     in.Seq,
		}

		_ = i.PreviousTxIDAddStr(*in.E.H)
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
				} else {
					fmt.Printf("no maatch %+v", cell)
				}
			}
		}

		lockingScript, _ := bscript.NewFromASM(strings.Join(lockScriptAsm, " "))
		o := &bt.Output{
			Satoshis:      *out.E.V,
			LockingScript: lockingScript,
		}

		tx.AddOutput(o)
	}

	return tx, nil
}
