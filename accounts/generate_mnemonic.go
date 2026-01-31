package accounts

import (
	"crypto/rand"
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

// MnemonicStrength represents the entropy strength for mnemonic generation.
type MnemonicStrength int

const (
	// Mnemonic128 generates a 12-word mnemonic (128 bits of entropy).
	Mnemonic128 MnemonicStrength = 128
	// Mnemonic160 generates a 15-word mnemonic (160 bits of entropy).
	Mnemonic160 MnemonicStrength = 160
	// Mnemonic192 generates an 18-word mnemonic (192 bits of entropy).
	Mnemonic192 MnemonicStrength = 192
	// Mnemonic224 generates a 21-word mnemonic (224 bits of entropy).
	Mnemonic224 MnemonicStrength = 224
	// Mnemonic256 generates a 24-word mnemonic (256 bits of entropy).
	Mnemonic256 MnemonicStrength = 256
)

// GenerateMnemonicOptions contains options for mnemonic generation.
type GenerateMnemonicOptions struct {
	// Strength is the entropy strength (128, 160, 192, 224, or 256 bits).
	// Default is 128 (12 words).
	Strength MnemonicStrength
}

// GenerateMnemonic generates a random mnemonic phrase.
//
// The default strength is 128 bits, which produces a 12-word mnemonic.
// Use Mnemonic256 for a more secure 24-word mnemonic.
//
// Example:
//
//	// Generate a 12-word mnemonic (default)
//	mnemonic, err := GenerateMnemonic()
//
//	// Generate a 24-word mnemonic
//	mnemonic, err := GenerateMnemonic(GenerateMnemonicOptions{Strength: Mnemonic256})
func GenerateMnemonic(opts ...GenerateMnemonicOptions) (string, error) {
	strength := Mnemonic128
	if len(opts) > 0 && opts[0].Strength > 0 {
		strength = opts[0].Strength
	}

	// Validate strength
	if strength != 128 && strength != 160 && strength != 192 && strength != 224 && strength != 256 {
		return "", fmt.Errorf("%w: strength must be 128, 160, 192, 224, or 256", ErrInvalidMnemonic)
	}

	// Generate entropy
	entropyBytes := int(strength) / 8
	entropy := make([]byte, entropyBytes)
	if _, err := rand.Read(entropy); err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Generate mnemonic from entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// MustGenerateMnemonic generates a mnemonic or panics on error.
func MustGenerateMnemonic(opts ...GenerateMnemonicOptions) string {
	mnemonic, err := GenerateMnemonic(opts...)
	if err != nil {
		panic(err)
	}
	return mnemonic
}

// ValidateMnemonic validates a mnemonic phrase.
//
// Example:
//
//	valid := ValidateMnemonic("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about")
//	// valid = true
func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// MnemonicToEntropy converts a mnemonic back to its entropy bytes.
func MnemonicToEntropy(mnemonic string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, ErrInvalidMnemonic
	}
	return bip39.EntropyFromMnemonic(mnemonic)
}

// EntropyToMnemonic converts entropy bytes to a mnemonic.
func EntropyToMnemonic(entropy []byte) (string, error) {
	return bip39.NewMnemonic(entropy)
}
