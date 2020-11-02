package bob

// InputAddresses returns the Bitcoin addresses for the transaction inputs
func (b *Tx) InputAddresses() (addresses []string) {
	for _, i := range b.In {
		if i.E.A != "false" {
			addresses = append(addresses, i.E.A)
		}
	}
	return
}

// OutputAddresses returns the Bitcoin addresses for the transaction outputs
func (b *Tx) OutputAddresses() (addresses []string) {
	for _, i := range b.In {
		if i.E.A != "false" {
			addresses = append(addresses, i.E.A)
		}
	}
	return
}
