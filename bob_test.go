package bob

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"testing"

	test "github.com/bitcoinschema/go-bob/testing"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

var sampleBobTx, sampleBobTxBadStrings, rawBobTx, parityBob, parityTx, boostTx, bigOrdTx string

func init() {
	sampleBobTx = test.GetTestHex("./testing/bob/207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d.json")
	sampleBobTxBadStrings = test.GetTestHex("./testing/bob/26b754e6fdf04121b8d91160a0b252a22ae30204fc552605b7f6d3f08419f29e.json")
	rawBobTx = test.GetTestHex("./testing/tx/2.hex")
	parityBob = test.GetTestHex("./testing/bob/98a5f6ef18eaea188bdfdc048f89a48af82627a15a76fd53584975f28ab3cc39.json")
	parityTx = test.GetTestHex("./testing/tx/98a5f6ef18eaea188bdfdc048f89a48af82627a15a76fd53584975f28ab3cc39.hex")
	boostTx = test.GetTestHex("./testing/tx/c5c7248302683107aa91014fd955908a7c572296e803512e497ddf7d1f458bd3.hex")
	bigOrdTx = test.GetTestHex("./testing/tx/c8cd6ff398d23e12e65ab065757fe6caf2d74b5e214b638365d61583030aa069.hex")
}

// TestNewFromBytes tests for nil case in NewFromBytes()
func TestNewFromBytes(t *testing.T) {

	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputLine        []byte
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				[]byte(""),
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				[]byte("invalid-json"),
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				[]byte(sampleBobTx),
				"0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000",
				"207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d",
				false,
				false,
			},
		}
	)

	// Run tests
	var b *Tx
	var err error

	for _, theTest := range tests {
		if b, err = NewFromBytes(theTest.inputLine); err != nil && !theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), theTest.inputLine, err.Error())
		} else if err == nil && theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), theTest.inputLine)
		} else if b == nil && !theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), theTest.inputLine)
		} else if b != nil && theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), theTest.inputLine)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			require.NoError(t, err)
			require.Equal(t, theTest.expectedTxString, str)
			require.Equal(t, theTest.expectedTxHash, b.Tx.Tx.H)
		}
	}
}

// ExampleNewFromBytes example using NewFromBytes()
func ExampleNewFromBytes() {
	b, err := NewFromBytes([]byte(sampleBobTx))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.Tx.H)
	// Output:found tx: 207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d
}

// BenchmarkNewFromBytes benchmarks the method NewFromBytes()
func BenchmarkNewFromBytes(b *testing.B) {
	tx := []byte(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = NewFromBytes(tx)
	}
}

// TestNewFromBytesPanic tests for nil case in NewFromBytes()
func TestNewFromBytesPanic(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		b, err := NewFromBytes([]byte(sampleBobTxBadStrings))
		require.NoError(t, err)
		require.NotNil(t, b)
		_, _ = b.ToRawTxString()
	})
}

// TestNewFromString tests for nil case in NewFromString()
func TestNewFromString(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputLine        string
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				"",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				"invalid-json",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				sampleBobTx,
				"0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000",
				"207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d",
				false,
				false,
			},
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, theTest := range tests {
		if b, err = NewFromString(theTest.inputLine); err != nil && !theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), theTest.inputLine, err.Error())
		} else if err == nil && theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), theTest.inputLine)
		} else if b == nil && !theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), theTest.inputLine)
		} else if b != nil && theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), theTest.inputLine)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			require.NoError(t, err)
			require.Equal(t, theTest.expectedTxString, str)
			require.Equal(t, theTest.expectedTxHash, b.Tx.Tx.H)
		}
	}
}

// ExampleNewFromString example using NewFromString()
func ExampleNewFromString() {
	b, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.Tx.H)
	// Output:found tx: 207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d
}

// BenchmarkNewFromString benchmarks the method NewFromString()
func BenchmarkNewFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewFromString(sampleBobTx)
	}
}

// TestNewFromRawTxString tests for nil case in NewFromRawTxString()
func TestNewFromRawTxString(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputLine        string
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				"",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				"0",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				"invalid-tx",
				"01000000000000000000",
				"",
				false,
				true,
			},
			// This does not work anymore with the new version of go-bt
			/*{
				rawBobTx,
				"01000000018f81a0884a11452aa5860f3b0016db1ec58d0cd654b2fa11ebdfd7e87eabeb0e00000000964c948f81a0884a11452aa5860f3b0016db1ec58d0cd654b2fa11ebdfd7e87eabeb0e020000006b483045022100bfbaa9cb07155cd3690722a9d527c70f91a6fc79233b0d091729e457e7c59dd902203059e1f077593654d8f7d2e22a5a40013e8dbf6fcccc5595305144149e5ed9014121039c555f098562d5f6cff2764008d6491961ab51c49356fee349720781ff6dfff7ffffffff00000000030000000000000000fd9c04006a2231394878696756345179427633744870515663554551797131707a5a56646f4175740a746578742f706c61696e04746578740a7477657463682e7478747c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b7477646174615f6a736f6e4dbd027b22637265617465645f6174223a22576564204f63742032312031323a30363a3238202b303030302032303230222c227477745f6964223a2231333138383836333639363530303033393639222c2274657874223a2257534a20456469746f7269616c20426f6172643a204a6f6520426964656e204d75737420416e73776572205175657374696f6e732041626f75742048756e74657220426964656e20616e64204368696e612068747470733a2f2f7777772e6272656974626172742e636f6d2f6e6174696f6e616c2d73656375726974792f323032302f31302f32302f77736a2d656469746f7269616c2d626f6172642d6a6f652d626964656e2d6d7573742d616e737765722d7175657374696f6e732d61626f75742d68756e7465722d626964656e2d616e642d6368696e612f2076696120404272656974626172744e657773204a6f6520426964656e206973206120746f74616c6c7920636f727275707420706f6c6974696369616e2c20616e6420676f74206361756768742e204174206c65617374206e6f7720686520776f6ee28099742062652061626c6520746f20726169736520796f7572205461786573202d204269676765737420696e63726561736520696e20552e532e20686973746f727921222c2275736572223a7b226e616d65223a22446f6e616c64204a2e205472756d70222c2273637265656e5f6e616d65223a227265616c446f6e616c645472756d70222c22637265617465645f6174223a22576564204d61722031382031333a34363a3338202b303030302032303039222c227477745f6964223a223235303733383737222c2270726f66696c655f696d6167655f75726c223a22687474703a2f2f7062732e7477696d672e636f6d2f70726f66696c655f696d616765732f3837343237363139373335373539363637322f6b5575687430306d5f6e6f726d616c2e6a7067227d7d0375726c3e68747470733a2f2f747769747465722e636f6d2f7265616c446f6e616c645472756d702f7374617475732f3133313838383633363936353030303339363907636f6d6d656e74046e756c6c076d625f75736572046e756c6c057265706c79046e756c6c047479706504706f73740974696d657374616d70046e756c6c036170700674776574636807696e766f6963652434626130313735632d313738662d346636332d623737662d3536323737313562326563657c22313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f454344534122313438574448366e465776356748383177657043726b3566486b4a774550415134514c58494531786378574a6b4e364a6538683361426d644161574947487841773333556167515951586539704672794b4a55334f786875324c54646b784b364d4b5675624a4475592f516957743164776f7a782b796167696c553d00000000000000001976a91405186ff0710ed004229e644c0653b2985c648a2388ac00000000000000001976a9142f0fadb49432be5f3d13a7db410e7c2ddae5103188ac00000000",
				"9ec47d91ff11edb62f337dc828c52e39072d1a5a2f1b180bbfae9c3279d81a7c",
				false,
				false,
			},
			{
				"01000000000100000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
				"01000000000100000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000",
				"c83a93b359a490c911b28b443f650ecc4da46def30b576249c47da5b610e2edc",
				false,
				false,
			},
			*/
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, theTest := range tests {
		if b, err = NewFromRawTxString(theTest.inputLine); err != nil && !theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), theTest.inputLine, err.Error())
		} else if err == nil && theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), theTest.inputLine)
		} else if b == nil && !theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), theTest.inputLine)
		} else if b != nil && theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), theTest.inputLine)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			require.NoError(t, err)
			require.Equal(t, theTest.expectedTxString, str)
			require.Equal(t, theTest.expectedTxHash, b.Tx.Tx.H)
		}
	}
}

// X5R

func TestMapFromRawTxString2(t *testing.T) {
	bob, err := NewFromRawTxString(`01000000018952fe8892c429e69feb9b2dd9cd1f12ed757dc62e8d628b5a215f78ed895374020000006a47304402204784632fabca0f4aaa05dd6983633b2e8bf708d8766d0385f3393fff0623b88c02201a760e144116d47967501c2ea50dc231ae57c0eb78769d63713ce0648025c820412103221cb24c4e8b05a58bcf2ee8411f62e337c8099c8646babd47d0960899f69acaffffffff04680b0000000000001976a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac0f720000000000001976a9146b1fe7b2063aa07766c764c0796fd4efd00340f288ac8a893b00000000001976a914be5f62df829ef754b8be09b37b04c4e7f9ff59d588ac0000000000000000ad006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f696403383038106f666665725f73657373696f6e5f6964403464303537386561643432393266653163643163393936643931623534613130653333653334623031396231386330613564353730376461346461346437653900000000`)

	if err != nil {
		t.Fatalf("error occurred: %s", err)
	} else if *bob.Out[3].Tape[1].Cell[5].S != "offer_click" {
		t.Fatalf("SET Failed %v", bob.Out[3].Tape[1].Cell[5])
	}
}

// 0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 534554 617070 746f6e6963706f77 74797065 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 383038 6f666665725f73657373696f6e5f6964 34643035373865616434323932666531636431633939366439316235346131306533336533346230313962313863306135643537303764613464613464376539

// ExampleNewFromRawTxString example using NewFromRawTxString()
func ExampleNewFromRawTxString() {
	b, err := NewFromRawTxString(rawBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.Tx.H)
	// Output:found tx: 9ec47d91ff11edb62f337dc828c52e39072d1a5a2f1b180bbfae9c3279d81a7c
}

// BenchmarkNewFromRawTxString benchmarks the method NewFromRawTxString()
func BenchmarkNewFromRawTxString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewFromRawTxString(rawBobTx)
	}
}

func testExampleTx() *transaction.Transaction {
	tx := transaction.NewTransaction()
	s := script.NewFromBytes([]byte{})
	_ = s.AppendOpcodes(script.OpFALSE, script.OpRETURN)
	_ = s.AppendPushDataArray([][]byte{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}, []byte{0x7c}, []byte("prefix2"), []byte("example data 2")})

	tx.AddOutput(&transaction.TransactionOutput{
		LockingScript: s,
		Satoshis:      0,
	})
	return tx
}

// TestNewFromTx tests for nil case in NewFromTx()
func TestNewFromTx(t *testing.T) {
	t.Parallel()

	validTx := testExampleTx()
	require.NotNil(t, validTx)

	var (
		// Testing private methods
		tests = []struct {
			inputTx          *transaction.Transaction
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				&transaction.Transaction{},
				"01000000000000000000",
				"f702453dd03b0f055e5437d76128141803984fb10acb85fc3b2184fae2f3fa78",
				false,
				false,
			},
			// This does not work anymore in the new version of go-bt
			/*{
				validTx,
				"010000000001000000000000000032006a07707265666978310c6578616d706c6520646174610213377c07707265666978320e6578616d706c652064617461203200000000",
				"f94e4adeac0cee5e9ff9985373622db9524e9f98d465dc024f85aec8acfeaf16",
				false,
				false,
			},*/
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, theTest := range tests {
		if b, err = NewFromTx(theTest.inputTx); err != nil && !theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), theTest.inputTx, err.Error())
		} else if err == nil && theTest.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), theTest.inputTx)
		} else if b == nil && !theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), theTest.inputTx)
		} else if b != nil && theTest.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), theTest.inputTx)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			require.NoError(t, err)
			require.Equal(t, theTest.expectedTxString, str)
			require.Equal(t, theTest.expectedTxHash, b.Tx.Tx.H)
		}
	}
}

func TestNewFromTxString(t *testing.T) {
	// BAP attestation
	testHex := test.GetTestHex("./testing/tx/98a5f6ef18eaea188bdfdc048f89a48af82627a15a76fd53584975f28ab3cc39.hex")

	bobTxFromString, err := NewFromRawTxString(testHex)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	bapOut := &bobTxFromString.Out[0].Tape[1]
	const bapPrefix = "1BAPSuaPnfGnSBM3GLV9yhxUdYe4vGbdMT"

	if *bapOut.Cell[0].S != bapPrefix {
		t.Errorf("Expected string(%s) is not same as actual string (%s)", bapPrefix, *bapOut.Cell[0].S)
	}

}

// ExampleNewFromTx example using NewFromTx()
func ExampleNewFromTx() {
	// Use an example TX
	exampleTx := testExampleTx()

	var b *Tx
	var err error
	if b, err = NewFromTx(exampleTx); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.Tx.H)
	// Output:found tx: f94e4adeac0cee5e9ff9985373622db9524e9f98d465dc024f85aec8acfeaf16
}

// BenchmarkNewFromTx benchmarks the method NewFromTx()
func BenchmarkNewFromTx(b *testing.B) {
	exampleTx := testExampleTx()
	for i := 0; i < b.N; i++ {
		_, _ = NewFromTx(exampleTx)
	}
}

// TestNewFromTx2 tests for nil case in NewFromTx()
func TestNewFromTx2(t *testing.T) {
	t.Parallel()
	b, err := NewFromTx(nil)
	require.Error(t, err)
	require.Nil(t, b)
}

// TestTx_ToTx tests for nil case in ToTx()
func TestTx_ToTx(t *testing.T) {

	bobTx, err := NewFromString(sampleBobTx)

	require.NoError(t, err)
	require.NotNil(t, bobTx)

	var tx *transaction.Transaction
	tx, err = bobTx.ToTx()
	require.NoError(t, err)
	require.NotNil(t, tx)

	// check that they have same number of ins and outs

	require.Equal(t, len(bobTx.In), len(tx.Inputs))
	require.Equal(t, len(bobTx.Out), len(tx.Outputs))

	parts, err := script.DecodeScript(*tx.Inputs[0].UnlockingScript)
	require.NoError(t, err)
	part0 := hex.EncodeToString(parts[0].Data)
	part1 := hex.EncodeToString(parts[1].Data)
	require.Equal(t, *bobTx.In[0].Tape[0].Cell[0].H, part0)
	require.Equal(t, *bobTx.In[0].Tape[0].Cell[1].H, part1)

	outParts, err := script.DecodeScript(*tx.Outputs[0].LockingScript)
	require.NoError(t, err)
	outPart1 := hex.EncodeToString(outParts[1].Data)

	log.Printf("%x ", outPart1)

	require.Equal(t, *bobTx.Out[0].Tape[0].Cell[0].Op, outParts[0].Op)

	require.Equal(t, bobTx.Tx.Tx.H, tx.TxID().String())
}

// ExampleTx_ToTx example using ToTx()
func ExampleTx_ToTx() {
	// Use an example TX
	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var tx *transaction.Transaction
	if tx, err = bobTx.ToTx(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("found tx: %s", tx.TxID())
	// Output:found tx: 207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d
}

// BenchmarkTx_ToTx benchmarks the method ToTx()
func BenchmarkTx_ToTx(b *testing.B) {
	bobTx, _ := NewFromString(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = bobTx.ToTx()
	}
}

// TestTx_ToRawTxString tests for nil case in ToRawTxString()
func TestTx_ToRawTxString(t *testing.T) {
	bobTx, err := NewFromString(sampleBobTx)
	require.NoError(t, err)
	require.NotNil(t, bobTx)

	testTx := "0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000"

	var rawTx string
	rawTx, err = bobTx.ToRawTxString()
	require.NoError(t, err)
	require.Equal(t, testTx, rawTx)
}

// ExampleTx_ToRawTxString example using ToRawTxString()
func ExampleTx_ToRawTxString() {
	// Use an example TX
	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var rawTx string
	if rawTx, err = bobTx.ToRawTxString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("found raw tx: %s", rawTx)
	// Output:found raw tx: 0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000
}

// BenchmarkTx_ToRawTxString benchmarks the method ToRawTxString()
func BenchmarkTx_ToRawTxString(b *testing.B) {
	bobTx, _ := NewFromString(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = bobTx.ToRawTxString()
	}
}

// Test parity with bmapjs
func Test_GoBT_ASM(t *testing.T) {

	t.Run("test go-bt ASM", func(t *testing.T) {

		// goBobTx, err := NewFromRawTxString(parityTx)
		// require.NotNil(t, goBobTx)
		// require.Nil(t, err)

		btTx, err := transaction.NewTransactionFromHex(parityTx)
		require.NoError(t, err)
		require.NotNil(t, *btTx)

		asmBt := btTx.Outputs[0].LockingScript.ToASM()

		pushDatas := strings.Split(asmBt, " ")
		require.Len(t, pushDatas, 10)
		require.Equal(t, "OP_RETURN", pushDatas[0])
	})

}

// Test parity with bmapjs
func TestBob_Vs_Bob(t *testing.T) {
	bmapjsTx, _ := NewFromString(parityBob)
	// import a tx from hex
	// TODO: - should this even work from a Bob string like this?
	goBobTx, _ := NewFromRawTxString(parityTx)

	// get same tx from bob output from bmapjs

	// make sure number of overall keys,  ins, outs, tapes, and cells are identical

	require.Equal(t, len(bmapjsTx.Out), len(goBobTx.Out))
	// require.Equal(t, len(bmapjsTx.Out[0].Tape), len(goBobTx.Out[0].Tape))
	// require.Equal(t, len(bmapjsTx.Out[0].Tape[1].Cell), len(goBobTx.Out[0].Tape[1].Cell))
	require.Equal(t, *bmapjsTx.Out[0].Tape[1].Cell[3].H, *goBobTx.Out[0].Tape[1].Cell[3].H)

	// fmt.Println(fmt.Sprintf("expected %+v", bmapjsTx.Out[0].Tape[1].Cell))
	// fmt.Println(fmt.Sprintf("actual %+v", goBobTx.Out[0].Tape[1].Cell))

	// require.Equal(t, len(bmapjsTx.Out[1].Tape), len(goBobTx.Out[1].Tape))
	// require.Equal(t, len(bmapjsTx.Out[1].Tape[0].Cell), len(goBobTx.Out[1].Tape[0].Cell))
	require.Equal(t, bmapjsTx.Tx.Tx.H, goBobTx.Tx.Tx.H)
}

// TestTx_ToString tests for nil case in ToString()
func TestTx_ToString(t *testing.T) {

	bobTx := new(Tx)
	err := bobTx.FromRawTxString(rawBobTx)
	require.NoError(t, err)

	// to string
	var txString string
	txString, err = bobTx.ToString()
	require.NoError(t, err)

	// make another bob tx from string
	var otherBob *Tx
	otherBob, err = NewFromString(txString)
	require.NoError(t, err)

	// check txid match
	require.Equal(t, bobTx.Tx.Tx.H, otherBob.Tx.Tx.H)
}

// TestTx_ToString2 example using ToString()
func TestTx_ToString2(t *testing.T) {
	// import a tx from hex
	// TODO: - should this even work from a Bob string like this?
	goBobTx, err := NewFromRawTxString(parityTx)
	require.NoError(t, err)

	if _, err = goBobTx.ToString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	}
	require.NoError(t, err)

}

// Test Boost
func TestTx_Boost(t *testing.T) {
	// import a tx from hex
	goBobTx, err := NewFromRawTxString(boostTx)
	require.NoError(t, err)

	if _, err = goBobTx.ToString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	}
	require.NoError(t, err)

	require.Len(t, goBobTx.Out[0].Tape[0].Cell, 89)
}

// Test HugeOrd
func TestRawTxString_HugeOrd(t *testing.T) {
	hugeOrdTx, err := transaction.NewTransactionFromHex(bigOrdTx)
	require.NoError(t, err)

	// import a tx from hex
	goBobTx, err := NewFromTx(hugeOrdTx)
	require.NoError(t, err)

	if _, err = goBobTx.ToString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	}
	require.NoError(t, err)

	require.Len(t, goBobTx.Out[0].Tape[0].Cell, 14)
	require.Len(t, goBobTx.Out[0].Tape, 2)
	require.Len(t, goBobTx.Out, 1437)

}

// Test HugeOrd From Tx
func TestTx_HugeOrd(t *testing.T) {
	// import a tx from hex
	goBobTx, err := NewFromRawTxString(bigOrdTx)
	require.NoError(t, err)

	if _, err = goBobTx.ToString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	}
	require.NoError(t, err)

	require.Len(t, goBobTx.Out[0].Tape[0].Cell, 14)
	require.Len(t, goBobTx.Out[0].Tape, 2)
	require.Len(t, goBobTx.Out, 1437)

}

// c5c7248302683107aa91014fd955908a7c572296e803512e497ddf7d1f458bd3
// BenchmarkTx_ToString benchmarks the method ToString()
func BenchmarkTx_ToString(b *testing.B) {
	bobTx, _ := NewFromString(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = bobTx.ToString()
	}
}
