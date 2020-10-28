## go-bob

Library for working with [BOB](https://bob.planaria.network/) formatted transactions, written in Go.

## Usage

```go
// Split BOB formatted NDJSON stream from Bitbus by line
line, err := reader.ReadBytes('\n')
if err == io.EOF {
  break
}

bobTx, err = bobData.NewFromBytes(line)
if err != nil {
  return err
}

```

bobTx will be of this type:

```go
type BobTx struct {
	ID  string   `json:"_id"`
	Blk Blk      `json:"blk"`
	Tx  TxInfo   `json:"tx"`
	In  []Input  `json:"in"`
	Out []Output `json:"out"`
}
```

## Helpers

### BOB From bytes

```go
// BOB from libsv.transaction
bobTx, err = NewFromTx(tx)

```

### BOB From libsv.Transaction

```go
// BOB from libsv.transaction
bobTx, err = NewFromTx(tx)

```

### BOB from raw tx string

```go
// BOB from raw tx string
bobTx, err := bob.NewFromRawTxString(rawTxString)

```

### BOB to libsv.Transaction

```go
// BOB from libsv.transaction
tx, err = bobTx.ToTx()

```

### BOB to string

```go
// BOB formatted JSON string
tx, err = bobTx.ToString()

```

### BOB to raw tx string

```go
// BOB to raw tx string
tx, err = bobTx.ToRawTxString()

```
