package signature

// VerifyMessage verifies that a message was signed by the provided address.
//
// Note: Only supports Externally Owned Accounts. Does not support Contract Accounts.
//
// Example:
//
//	valid, err := VerifyMessage(
//		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
//		NewSignableMessage("hello world"),
//		"0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c",
//	)
//	// true
func VerifyMessage(address string, message SignableMessage, signature any) (bool, error) {
	recoveredAddress, err := RecoverMessageAddress(message, signature)
	if err != nil {
		return false, err
	}

	return isAddressEqual(address, recoveredAddress), nil
}
