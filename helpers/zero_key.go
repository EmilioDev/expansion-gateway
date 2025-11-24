package helpers

// Zero the private key material (call when you're done with the private scalar)
func ZeroKey(k *[32]byte) {
	for i := range k {
		k[i] = 0
	}
}
