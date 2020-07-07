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
