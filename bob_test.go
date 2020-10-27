package bob

import (
	"testing"

	"github.com/bitcoinschema/go-bitcoin"
)

const sampleBobTx = `{ "_id": "5ed082db57cd6b1658b88400", "tx": { "h": "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" }, "in": [ { "i": 0, "tape": [ { "cell": [ { "b": "MEUCIQDwEsO9N4EJGqjlPKsv/LkKzO2MZVALQQhv0iXkjJjB1wIgC4/xF7js0rLX6VVRvFobO7zKgEmGRHmii+2dyEKoaARB", "s": "0E\u0002!\u0000�\u0012ý7�\t\u001a��<�/��\n��eP\u000bA\bo�%䌘��\u0002 \u000b��\u0017��Ҳ��UQ�Z\u001b;�ʀI�Dy����B�h\u0004A", "ii": 0, "i": 0 }, { "b": "A+9bsilk1SnAr3SNmmOBQy8FKY56Zu0v4i55dbFQJSin", "s": "\u0003�[�)d�)��t��c�C/\u0005)�zf�/�.yu�P%(�", "ii": 1, "i": 1 } ], "i": 0 } ], "e": { "h": "3d1fc854830cb7f5cf4e89459f1e2f4331ffed09ad66a02ce1140c553c9d5af1", "i": 1, "a": "1FFuYLM8a66GddCG25nUbarazeMr5dnUwC" }, "seq": 4294967295 } ], "out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "ii": 0, "i": 0 }, { "op": 106, "ops": "OP_RETURN", "ii": 1, "i": 1 } ], "i": 0 }, { "cell": [ { "b": "5LiA54Gv6IO96Zmk5Y2D5bm05pqX", "s": "一灯能除千年暗", "ii": 2, "i": 0 }, { "b": "NThhNTk3", "s": "58a597", "ii": 3, "i": 1 } ], "i": 1 } ], "e": { "v": 0, "i": 0, "a": "false" } }, { "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "ii": 0, "i": 0 }, { "op": 169, "ops": "OP_HASH160", "ii": 1, "i": 1 }, { "b": "nGNxXG0fpsYbMdKRFRbhw9s736g=", "s": "�cq\\m\u001f��\u001b1ґ\u0015\u0016���;ߨ", "ii": 2, "i": 2 }, { "op": 136, "ops": "OP_EQUALVERIFY", "ii": 3, "i": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "ii": 4, "i": 4 } ], "i": 0 } ], "e": { "v": 111411, "i": 1, "a": "1FFuYLM8a66GddCG25nUbarazeMr5dnUwC" } } ], "lock": 0, "blk": { "i": 635140, "h": "0000000000000000031d01ce0a8471d6cfab81d403ba10c878f671eac28d5d39", "t": 1589607858 }, "i": 4042 }`

const sampleBobTxBadStrings = `{ "_id": "5f08ddeed1352a2c3432f4db", "tx": { "h": "26b754e6fdf04121b8d91160a0b252a22ae30204fc552605b7f6d3f08419f29e" }, "in": [ { "i": 0, "seq": 4294967295, "tape": [ { "cell": [ { "s": "0E\u0002!\u0000����;�Z��\b\th�&���5����6��` + "`" + `\u0016�Z�N\u0002 WUI\u001bz)\nE{\u001f��0�g�꨻*}\u0018QV��dO�D@�A", "h": "3045022100afbbffff3bb55aaec20809689026acbccf35bcb4e2f29c36aaf86016d85abe4e02205755491b7a290a457b1fbea2308567ddeaa8bb2a7d185156a1f3644f854440d941", "b": "MEUCIQCvu///O7VarsIICWiQJqy8zzW8tOLynDaq+GAW2Fq+TgIgV1VJG3opCkV7H76iMIVn3eqouyp9GFFWofNkT4VEQNlB", "i": 0, "ii": 0 }, { "s": "\u0004@��8��x��x���,#\u001d�(��B�A%\f����E��\u0000��T[�=(�\u0017Ϳ\u0001\u0010*\u001cr\\iZ��\u0007Ha�\u0018WM�(", "h": "0440ffb338848f78bfbb78b9b4a82c231dc728ceef42b341250c84ba99cf458bf2af0095df545bef3d28e717cdbf01102a1c725c695adfe40748619518574df228", "b": "BED/sziEj3i/u3i5tKgsIx3HKM7vQrNBJQyEupnPRYvyrwCV31Rb7z0o5xfNvwEQKhxyXGla3+QHSGGVGFdN8ig=", "i": 1, "ii": 1 } ], "i": 0 } ], "e": { "h": "744a55a8637aa191aa058630da51803abbeadc2de3d65b4acace1f5f10789c5b", "i": 1, "a": "1LC16EQVsqVYGeYTCrjvNf8j28zr4DwBuk" } } ], "out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "i": 0, "ii": 0 }, { "op": 106, "ops": "OP_RETURN", "i": 1, "ii": 1 } ], "i": 0 }, { "cell": [ { "s": "1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT", "h": "31424150537561506e66476e53424d33474c56397968785564596534764762644d54", "b": "MUJBUFN1YVBuZkduU0JNM0dMVjl5aHhVZFllNHZHYmRNVA==", "i": 0, "ii": 2 }, { "s": "ATTEST", "h": "415454455354", "b": "QVRURVNU", "i": 1, "ii": 3 }, { "s": "16ca90ce3c6347132adba40aa0d5faa3b2bf2015678ffc63db1511b676885e25", "h": "31366361393063653363363334373133326164626134306161306435666161336232626632303135363738666663363364623135313162363736383835653235", "b": "MTZjYTkwY2UzYzYzNDcxMzJhZGJhNDBhYTBkNWZhYTNiMmJmMjAxNTY3OGZmYzYzZGIxNTExYjY3Njg4NWUyNQ==", "i": 2, "ii": 4 }, { "s": "0", "h": "30", "b": "MA==", "i": 3, "ii": 5 } ], "i": 1 }, { "cell": [ { "s": "15PciHG22SNLQJXMoSUaWVi7WSqc7hCfva", "h": "313550636948473232534e4c514a584d6f5355615756693757537163376843667661", "b": "MTVQY2lIRzIyU05MUUpYTW9TVWFXVmk3V1NxYzdoQ2Z2YQ==", "i": 0, "ii": 7 }, { "s": "BITCOIN_ECDSA", "h": "424954434f494e5f4543445341", "b": "QklUQ09JTl9FQ0RTQQ==", "i": 1, "ii": 8 }, { "s": "134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da", "h": "31333461365458787a675139417a33773842637667645a7941355571524c38396461", "b": "MTM0YTZUWHh6Z1E5QXozdzhCY3ZnZFp5QTVVcVJMODlkYQ==", "i": 2, "ii": 9 }, { "s": "\u001f�V���j{k�\u0010ҕ�QA�]�Ӛ` + "`" + `7N����^���)YΓ\u001f@�qWcH}�V��Y�\u0019F�C�V�@�\r�a�", "h": "1fc756c3fcc76a7b6bcf10d295a75141ef5dbbd39a60374ea796eb92d85e84a0a32959ce931f40dc715763487de7a856acca59fc19468343b4569340d20d9761ed", "b": "H8dWw/zHantrzxDSladRQe9du9OaYDdOp5brkthehKCjKVnOkx9A3HFXY0h956hWrMpZ/BlGg0O0VpNA0g2XYe0=", "i": 3, "ii": 10 } ], "i": 2 } ], "e": { "v": 0, "i": 0, "a": "false" } }, { "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "i": 0, "ii": 0 }, { "op": 169, "ops": "OP_HASH160", "i": 1, "ii": 1 }, { "s": "�\no;L˺��E\t^��{i\u0011}", "h": "d27f0a6f3b4ccbbacaf945095ed3eeb97b69117d", "b": "0n8KbztMy7rK+UUJXtPuuXtpEX0=", "i": 2, "ii": 2 }, { "op": 136, "ops": "OP_EQUALVERIFY", "i": 3, "ii": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "i": 4, "ii": 4 } ], "i": 0 } ], "e": { "v": 14491552, "i": 1, "a": "1LC16EQVsqVYGeYTCrjvNf8j28zr4DwBuk" } } ], "lock": 0, "timestamp": 1594416622135 }`

func TestFromString(t *testing.T) {
	bobData := New()
	err := bobData.FromString(sampleBobTx)
	if err != nil {
		t.Error(err)
	}

	if bobData.Tx.H != "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" {
		t.Error("From String Failed")
	}

}

func TestFromTx(t *testing.T) {
	pk := "80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01"

	privateKey, err := bitcoin.PrivateKeyFromString(pk)
	if err != nil {
		t.Errorf("Failed to get private key")
	}
	opReturn1 := bitcoin.OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}, []byte{0x7c}, []byte("prefix2"), []byte("example data 2")}
	tx, err := bitcoin.CreateTx(nil, nil, []bitcoin.OpReturnData{opReturn1}, privateKey)
	if err != nil {
		t.Errorf("Failed to create tx %s", err)
	}

	bobTx := New()
	err = bobTx.FromTx(tx)
	if err != nil {
		t.Errorf("Failed to create from bob tx %s", err)
	}

}

func TestFromBadString(t *testing.T) {
	bobBadStrings := New()
	err := bobBadStrings.FromString(sampleBobTxBadStrings)
	if err != nil {
		t.Error(err)
	}
	if bobBadStrings.Tx.H != "26b754e6fdf04121b8d91160a0b252a22ae30204fc552605b7f6d3f08419f29e" {
		t.Error("From String Failed")
	}
}

func TestFromBytes(t *testing.T) {
	bobData := New()
	bobData.FromBytes([]byte(sampleBobTx))

	if bobData.Tx.H != "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" {
		t.Error("From Bytes Failed")
	}
}
