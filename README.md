## go-bob

Library for working with [BOB](https://) formatted transactions, written in Go.

## Usage

```go
// Split BOB formatted NDJSON stream from Bitbus by line
line, err := reader.ReadBytes('\n')
if err == io.EOF {
  break
}

bobData := bob.New()
err = bobData.FromBytes(line)
if err != nil {
  return err
}

```

bobData will be of this type:

```go
type Tx struct {
	ID  string   `json:"_id"`
	Blk Blk      `json:"blk"`
	Tx  TxInfo   `json:"tx"`
	In  []Input  `json:"in"`
	Out []Output `json:"out"`
}
```

## Features

### BOB From libsv.Transaction

```go
// BOB from libsv.transaction
bobData := bob.New()
err = bobData.FromTx(tx)

```

### BOB to libsv.Transaction

```go
// BOB from libsv.transaction
bobData := bob.New()
tx, err = bobData.ToTx()

```

### BOB to string

```go
// BOB formatted JSON string
bobData := bob.New()
tx, err = bobData.ToString()

```

### BOB to raw tx string

```go
// BOB to raw tx string
bobData := bob.New()
tx, err = bobData.ToRawTxString()

```

## BOB from raw tx string

```go
// BOB from raw tx string
bobTx := New()
err := bobTx.FromRawTxString(rawTxString)

```
