package bob

import (
	"testing"

	"github.com/bitcoinschema/go-bitcoin"
)

const sampleBobTx = `{ "_id": "5ed082db57cd6b1658b88400", "tx": { "h": "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" }, "in": [ { "i": 0, "tape": [ { "cell": [ { "b": "MEUCIQDwEsO9N4EJGqjlPKsv/LkKzO2MZVALQQhv0iXkjJjB1wIgC4/xF7js0rLX6VVRvFobO7zKgEmGRHmii+2dyEKoaARB", "s": "0E\u0002!\u0000�\u0012ý7�\t\u001a��<�/��\n��eP\u000bA\bo�%䌘��\u0002 \u000b��\u0017��Ҳ��UQ�Z\u001b;�ʀI�Dy����B�h\u0004A", "ii": 0, "i": 0 }, { "b": "A+9bsilk1SnAr3SNmmOBQy8FKY56Zu0v4i55dbFQJSin", "s": "\u0003�[�)d�)��t��c�C/\u0005)�zf�/�.yu�P%(�", "ii": 1, "i": 1 } ], "i": 0 } ], "e": { "h": "3d1fc854830cb7f5cf4e89459f1e2f4331ffed09ad66a02ce1140c553c9d5af1", "i": 1, "a": "1FFuYLM8a66GddCG25nUbarazeMr5dnUwC" }, "seq": 4294967295 } ], "out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "ii": 0, "i": 0 }, { "op": 106, "ops": "OP_RETURN", "ii": 1, "i": 1 } ], "i": 0 }, { "cell": [ { "b": "5LiA54Gv6IO96Zmk5Y2D5bm05pqX", "s": "一灯能除千年暗", "ii": 2, "i": 0 }, { "b": "NThhNTk3", "s": "58a597", "ii": 3, "i": 1 } ], "i": 1 } ], "e": { "v": 0, "i": 0, "a": "false" } }, { "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "ii": 0, "i": 0 }, { "op": 169, "ops": "OP_HASH160", "ii": 1, "i": 1 }, { "b": "nGNxXG0fpsYbMdKRFRbhw9s736g=", "s": "�cq\\m\u001f��\u001b1ґ\u0015\u0016���;ߨ", "ii": 2, "i": 2 }, { "op": 136, "ops": "OP_EQUALVERIFY", "ii": 3, "i": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "ii": 4, "i": 4 } ], "i": 0 } ], "e": { "v": 111411, "i": 1, "a": "1FFuYLM8a66GddCG25nUbarazeMr5dnUwC" } } ], "lock": 0, "blk": { "i": 635140, "h": "0000000000000000031d01ce0a8471d6cfab81d403ba10c878f671eac28d5d39", "t": 1589607858 }, "i": 4042 }`
const sampleBobTxBadStrings = `{ "_id": "5f08ddeed1352a2c3432f4db", "tx": { "h": "26b754e6fdf04121b8d91160a0b252a22ae30204fc552605b7f6d3f08419f29e" }, "in": [ { "i": 0, "seq": 4294967295, "tape": [ { "cell": [ { "s": "0E\u0002!\u0000����;�Z��\b\th�&���5����6��` + "`" + `\u0016�Z�N\u0002 WUI\u001bz)\nE{\u001f��0�g�꨻*}\u0018QV��dO�D@�A", "h": "3045022100afbbffff3bb55aaec20809689026acbccf35bcb4e2f29c36aaf86016d85abe4e02205755491b7a290a457b1fbea2308567ddeaa8bb2a7d185156a1f3644f854440d941", "b": "MEUCIQCvu///O7VarsIICWiQJqy8zzW8tOLynDaq+GAW2Fq+TgIgV1VJG3opCkV7H76iMIVn3eqouyp9GFFWofNkT4VEQNlB", "i": 0, "ii": 0 }, { "s": "\u0004@��8��x��x���,#\u001d�(��B�A%\f����E��\u0000��T[�=(�\u0017Ϳ\u0001\u0010*\u001cr\\iZ��\u0007Ha�\u0018WM�(", "h": "0440ffb338848f78bfbb78b9b4a82c231dc728ceef42b341250c84ba99cf458bf2af0095df545bef3d28e717cdbf01102a1c725c695adfe40748619518574df228", "b": "BED/sziEj3i/u3i5tKgsIx3HKM7vQrNBJQyEupnPRYvyrwCV31Rb7z0o5xfNvwEQKhxyXGla3+QHSGGVGFdN8ig=", "i": 1, "ii": 1 } ], "i": 0 } ], "e": { "h": "744a55a8637aa191aa058630da51803abbeadc2de3d65b4acace1f5f10789c5b", "i": 1, "a": "1LC16EQVsqVYGeYTCrjvNf8j28zr4DwBuk" } } ], "out": [ { "i": 0, "tape": [ { "cell": [ { "op": 0, "ops": "OP_0", "i": 0, "ii": 0 }, { "op": 106, "ops": "OP_RETURN", "i": 1, "ii": 1 } ], "i": 0 }, { "cell": [ { "s": "1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT", "h": "31424150537561506e66476e53424d33474c56397968785564596534764762644d54", "b": "MUJBUFN1YVBuZkduU0JNM0dMVjl5aHhVZFllNHZHYmRNVA==", "i": 0, "ii": 2 }, { "s": "ATTEST", "h": "415454455354", "b": "QVRURVNU", "i": 1, "ii": 3 }, { "s": "16ca90ce3c6347132adba40aa0d5faa3b2bf2015678ffc63db1511b676885e25", "h": "31366361393063653363363334373133326164626134306161306435666161336232626632303135363738666663363364623135313162363736383835653235", "b": "MTZjYTkwY2UzYzYzNDcxMzJhZGJhNDBhYTBkNWZhYTNiMmJmMjAxNTY3OGZmYzYzZGIxNTExYjY3Njg4NWUyNQ==", "i": 2, "ii": 4 }, { "s": "0", "h": "30", "b": "MA==", "i": 3, "ii": 5 } ], "i": 1 }, { "cell": [ { "s": "15PciHG22SNLQJXMoSUaWVi7WSqc7hCfva", "h": "313550636948473232534e4c514a584d6f5355615756693757537163376843667661", "b": "MTVQY2lIRzIyU05MUUpYTW9TVWFXVmk3V1NxYzdoQ2Z2YQ==", "i": 0, "ii": 7 }, { "s": "BITCOIN_ECDSA", "h": "424954434f494e5f4543445341", "b": "QklUQ09JTl9FQ0RTQQ==", "i": 1, "ii": 8 }, { "s": "134a6TXxzgQ9Az3w8BcvgdZyA5UqRL89da", "h": "31333461365458787a675139417a33773842637667645a7941355571524c38396461", "b": "MTM0YTZUWHh6Z1E5QXozdzhCY3ZnZFp5QTVVcVJMODlkYQ==", "i": 2, "ii": 9 }, { "s": "\u001f�V���j{k�\u0010ҕ�QA�]�Ӛ` + "`" + `7N����^���)YΓ\u001f@�qWcH}�V��Y�\u0019F�C�V�@�\r�a�", "h": "1fc756c3fcc76a7b6bcf10d295a75141ef5dbbd39a60374ea796eb92d85e84a0a32959ce931f40dc715763487de7a856acca59fc19468343b4569340d20d9761ed", "b": "H8dWw/zHantrzxDSladRQe9du9OaYDdOp5brkthehKCjKVnOkx9A3HFXY0h956hWrMpZ/BlGg0O0VpNA0g2XYe0=", "i": 3, "ii": 10 } ], "i": 2 } ], "e": { "v": 0, "i": 0, "a": "false" } }, { "i": 1, "tape": [ { "cell": [ { "op": 118, "ops": "OP_DUP", "i": 0, "ii": 0 }, { "op": 169, "ops": "OP_HASH160", "i": 1, "ii": 1 }, { "s": "�\no;L˺��E\t^��{i\u0011}", "h": "d27f0a6f3b4ccbbacaf945095ed3eeb97b69117d", "b": "0n8KbztMy7rK+UUJXtPuuXtpEX0=", "i": 2, "ii": 2 }, { "op": 136, "ops": "OP_EQUALVERIFY", "i": 3, "ii": 3 }, { "op": 172, "ops": "OP_CHECKSIG", "i": 4, "ii": 4 } ], "i": 0 } ], "e": { "v": 14491552, "i": 1, "a": "1LC16EQVsqVYGeYTCrjvNf8j28zr4DwBuk" } } ], "lock": 0, "timestamp": 1594416622135 }`

func TestNewFromString(t *testing.T) {
	bobData, err := NewFromString(sampleBobTx)
	if err != nil {
		t.Error(err)
	}

	if bobData.Tx.H != "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" {
		t.Error("From String Failed", bobData.Tx.H)
	}

}

func TestNewFromTx(t *testing.T) {
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

	bobTx, err := NewFromTx(tx)
	if err != nil {
		t.Errorf("Failed to create from bob tx %s", err)
	}

	if bobTx.Tx.H != "f94e4adeac0cee5e9ff9985373622db9524e9f98d465dc024f85aec8acfeaf16" {
		t.Error("Tx hash doesnt match", bobTx.Tx.H)
	}

}

func TestToTx(t *testing.T) {

	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		t.Errorf("Failed to create bob tx %s", err)
	}

	tx, err := bobTx.ToTx()
	if err != nil {
		t.Errorf("Failed to create tx %s", err)
	}

	if tx.GetTxID() != bobTx.Tx.H {
		t.Errorf("Unexpected tx result %s %s %s", tx.ToString(), bobTx.Tx.H, tx.GetTxID())
	}
}

func TestToRawTxString(t *testing.T) {
	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		t.Errorf("Failed to create bob tx %s", err)
	}

	rawTx, err := bobTx.ToRawTxString()
	if err != nil || rawTx != "0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000" {
		t.Errorf("Failed to convert Bob tx to raw tx string %s %s", rawTx, err)
	}
}

func TestToString(t *testing.T) {
	rawTxString := "0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000"

	bobTx := new(Tx)
	err := bobTx.FromRawTxString(rawTxString)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	// to string
	txString, err := bobTx.ToString()
	if err != nil {
		t.Errorf("Error %s", err)
	}

	// make another bob tx from string
	otherBob, err := NewFromString(txString)
	if err != nil {
		t.Errorf("Failed to create bob tx %s", err)
	}

	// check txid match
	if bobTx.Tx.H != otherBob.Tx.H {
		t.Errorf("TXIDS do not match!")
	}
}

func TestFromRawTx(t *testing.T) {

	// rawTxString := "0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000"
	rawTxString := "01000000018f81a0884a11452aa5860f3b0016db1ec58d0cd654b2fa11ebdfd7e87eabeb0e020000006b483045022100bfbaa9cb07155cd3690722a9d527c70f91a6fc79233b0d091729e457e7c59dd902203059e1f077593654d8f7d2e22a5a40013e8dbf6fcccc5595305144149e5ed9014121039c555f098562d5f6cff2764008d6491961ab51c49356fee349720781ff6dfff7ffffffff030000000000000000fda004006a2231394878696756345179427633744870515663554551797131707a5a56646f41757401200a746578742f706c61696e04746578740a7477657463682e747874017c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b7477646174615f6a736f6e4dbd027b22637265617465645f6174223a22576564204f63742032312031323a30363a3238202b303030302032303230222c227477745f6964223a2231333138383836333639363530303033393639222c2274657874223a2257534a20456469746f7269616c20426f6172643a204a6f6520426964656e204d75737420416e73776572205175657374696f6e732041626f75742048756e74657220426964656e20616e64204368696e612068747470733a2f2f7777772e6272656974626172742e636f6d2f6e6174696f6e616c2d73656375726974792f323032302f31302f32302f77736a2d656469746f7269616c2d626f6172642d6a6f652d626964656e2d6d7573742d616e737765722d7175657374696f6e732d61626f75742d68756e7465722d626964656e2d616e642d6368696e612f2076696120404272656974626172744e657773204a6f6520426964656e206973206120746f74616c6c7920636f727275707420706f6c6974696369616e2c20616e6420676f74206361756768742e204174206c65617374206e6f7720686520776f6ee28099742062652061626c6520746f20726169736520796f7572205461786573202d204269676765737420696e63726561736520696e20552e532e20686973746f727921222c2275736572223a7b226e616d65223a22446f6e616c64204a2e205472756d70222c2273637265656e5f6e616d65223a227265616c446f6e616c645472756d70222c22637265617465645f6174223a22576564204d61722031382031333a34363a3338202b303030302032303039222c227477745f6964223a223235303733383737222c2270726f66696c655f696d6167655f75726c223a22687474703a2f2f7062732e7477696d672e636f6d2f70726f66696c655f696d616765732f3837343237363139373335373539363637322f6b5575687430306d5f6e6f726d616c2e6a7067227d7d0375726c3e68747470733a2f2f747769747465722e636f6d2f7265616c446f6e616c645472756d702f7374617475732f3133313838383633363936353030303339363907636f6d6d656e74046e756c6c076d625f75736572046e756c6c057265706c79046e756c6c047479706504706f73740974696d657374616d70046e756c6c036170700674776574636807696e766f6963652434626130313735632d313738662d346636332d623737662d353632373731356232656365017c22313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f454344534122313438574448366e465776356748383177657043726b3566486b4a774550415134514c58494531786378574a6b4e364a6538683361426d644161574947487841773333556167515951586539704672794b4a55334f786875324c54646b784b364d4b5675624a4475592f516957743164776f7a782b796167696c553deb100000000000001976a91405186ff0710ed004229e644c0653b2985c648a2388ace4350900000000001976a9142f0fadb49432be5f3d13a7db410e7c2ddae5103188ac00000000"
	bobTx := new(Tx)
	err := bobTx.FromRawTxString(rawTxString)
	if err != nil {
		t.Errorf("Error %s", err)
	}
}

func TestFromBadString(t *testing.T) {
	bobBadStrings, err := NewFromString(sampleBobTxBadStrings)
	if err != nil {
		t.Error(err)
	}
	if bobBadStrings.Tx.H != "26b754e6fdf04121b8d91160a0b252a22ae30204fc552605b7f6d3f08419f29e" {
		t.Error("From String Failed")
	}
}

func TestNewFromBytes(t *testing.T) {

	bobTx, err := NewFromBytes([]byte(sampleBobTx))
	if err != nil {
		t.Errorf("Failed to create bob tx %s", err)
	}

	if bobTx.Tx.H != "207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d" {
		t.Error("From Bytes Failed")
	}
}
