package accounts

import (
	"encoding/hex"
	"fmt"

	"github.com/tyler-smith/go-bip32"
)

// HDKeyWrapper wraps a bip32.Key to implement the HDKey interface.
type HDKeyWrapper struct {
	key *bip32.Key
}

// NewHDKeyWrapper creates a new HDKeyWrapper from a bip32.Key.
func NewHDKeyWrapper(key *bip32.Key) *HDKeyWrapper {
	return &HDKeyWrapper{key: key}
}

// Derive derives a child key at the given path.
func (h *HDKeyWrapper) Derive(path string) (HDKey, error) {
	// Parse the path and derive step by step
	derivedKey, err := deriveFromPath(h.key, path)
	if err != nil {
		return nil, err
	}
	return &HDKeyWrapper{key: derivedKey}, nil
}

// PrivateKey returns the private key bytes.
func (h *HDKeyWrapper) PrivateKey() []byte {
	return h.key.Key
}

// PublicKey returns the public key bytes.
func (h *HDKeyWrapper) PublicKey() []byte {
	return h.key.PublicKey().Key
}

// HDKeyToAccountOptions contains options for creating an HD account.
type HDKeyToAccountOptions struct {
	HDOptions
}

// HDKeyToAccount creates an Account from a HD Key.
//
// Example:
//
//	// Create from HD key with default path m/44'/60'/0'/0/0
//	account, err := HDKeyToAccount(hdKey)
//
//	// Create with custom account index
//	account, err := HDKeyToAccount(hdKey, HDKeyToAccountOptions{
//		HDOptions: HDOptions{AccountIndex: 1},
//	})
//
//	// Create with custom path
//	account, err := HDKeyToAccount(hdKey, HDKeyToAccountOptions{
//		HDOptions: HDOptions{Path: "m/44'/60'/0'/0/5"},
//	})
func HDKeyToAccount(hdKey HDKey, opts ...HDKeyToAccountOptions) (*HDAccount, error) {
	var options HDKeyToAccountOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// Determine the derivation path
	path := options.Path
	if path == "" {
		path = DefaultHDPath(options.AccountIndex, options.ChangeIndex, options.AddressIndex)
	}

	// Derive the key at the path
	derivedKey, err := hdKey.Derive(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidHDPath, err)
	}

	// Get the private key as hex
	privateKeyBytes := derivedKey.PrivateKey()
	privateKey := "0x" + hex.EncodeToString(privateKeyBytes)

	// Create a private key account
	pkAccount, err := PrivateKeyToAccount(privateKey)
	if err != nil {
		return nil, err
	}

	// Update the source to HD
	pkAccount.Source = AccountSourceHD

	return &HDAccount{
		LocalAccount: pkAccount.LocalAccount,
		hdKey:        derivedKey,
	}, nil
}

// MustHDKeyToAccount creates an HD account or panics on error.
func MustHDKeyToAccount(hdKey HDKey, opts ...HDKeyToAccountOptions) *HDAccount {
	account, err := HDKeyToAccount(hdKey, opts...)
	if err != nil {
		panic(err)
	}
	return account
}

// deriveFromPath derives a key from a BIP32 path string.
func deriveFromPath(masterKey *bip32.Key, path string) (*bip32.Key, error) {
	// Parse the path
	indices, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, index := range indices {
		var err error
		key, err = key.NewChildKey(index)
		if err != nil {
			return nil, fmt.Errorf("failed to derive child key: %w", err)
		}
	}

	return key, nil
}

// parsePath parses a BIP32 derivation path string.
// Supports paths like "m/44'/60'/0'/0/0"
func parsePath(path string) ([]uint32, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	// Remove "m/" prefix if present
	if len(path) >= 2 && path[0] == 'm' && path[1] == '/' {
		path = path[2:]
	}

	if len(path) == 0 {
		return []uint32{}, nil
	}

	var indices []uint32
	var current uint32
	var hasDigit bool
	hardened := false

	for i := 0; i <= len(path); i++ {
		var c byte
		if i < len(path) {
			c = path[i]
		}

		if i == len(path) || c == '/' {
			if !hasDigit {
				if i == len(path) && len(indices) == 0 {
					break
				}
				return nil, fmt.Errorf("invalid path segment")
			}

			if hardened {
				// Add hardened offset (2^31)
				current += 0x80000000
			}

			indices = append(indices, current)
			current = 0
			hasDigit = false
			hardened = false
		} else if c == '\'' || c == 'h' || c == 'H' {
			hardened = true
		} else if c >= '0' && c <= '9' {
			hasDigit = true
			current = current*10 + uint32(c-'0')
		} else {
			return nil, fmt.Errorf("invalid character in path: %c", c)
		}
	}

	return indices, nil
}
