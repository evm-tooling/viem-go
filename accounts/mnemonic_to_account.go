package accounts

import (
	"fmt"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// MnemonicToAccountOptions contains options for creating an account from a mnemonic.
type MnemonicToAccountOptions struct {
	HDOptions
	// Passphrase is an optional passphrase for the mnemonic (BIP39 passphrase).
	Passphrase string
}

// MnemonicToAccount creates an Account from a mnemonic phrase.
//
// The mnemonic is converted to a seed using BIP39, then to an HD key using BIP32,
// and finally to an account at the specified derivation path.
//
// Example:
//
//	// Create account with default path m/44'/60'/0'/0/0
//	account, err := MnemonicToAccount("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about")
//
//	// Create account with passphrase
//	account, err := MnemonicToAccount("abandon ...", MnemonicToAccountOptions{
//		Passphrase: "my-passphrase",
//	})
//
//	// Create account at different index
//	account, err := MnemonicToAccount("abandon ...", MnemonicToAccountOptions{
//		HDOptions: HDOptions{AccountIndex: 1},
//	})
func MnemonicToAccount(mnemonic string, opts ...MnemonicToAccountOptions) (*HDAccount, error) {
	var options MnemonicToAccountOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// Validate the mnemonic
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	// Convert mnemonic to seed (with optional passphrase)
	seed := bip39.NewSeed(mnemonic, options.Passphrase)

	// Create master key from seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	// Wrap the master key
	hdKey := NewHDKeyWrapper(masterKey)

	// Create HD account
	hdOpts := HDKeyToAccountOptions{
		HDOptions: options.HDOptions,
	}

	account, err := HDKeyToAccount(hdKey, hdOpts)
	if err != nil {
		return nil, err
	}

	// Update source to mnemonic
	account.Source = AccountSourceMnemonic

	return account, nil
}

// MustMnemonicToAccount creates an account from a mnemonic or panics on error.
func MustMnemonicToAccount(mnemonic string, opts ...MnemonicToAccountOptions) *HDAccount {
	account, err := MnemonicToAccount(mnemonic, opts...)
	if err != nil {
		panic(err)
	}
	return account
}

// MnemonicToSeed converts a mnemonic to a BIP39 seed.
//
// Example:
//
//	seed := MnemonicToSeed("abandon ...", "passphrase")
func MnemonicToSeed(mnemonic string, passphrase string) []byte {
	return bip39.NewSeed(mnemonic, passphrase)
}

// SeedToHDKey converts a BIP39 seed to an HD key.
//
// Example:
//
//	hdKey, err := SeedToHDKey(seed)
func SeedToHDKey(seed []byte) (HDKey, error) {
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}
	return NewHDKeyWrapper(masterKey), nil
}
