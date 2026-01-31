package accounts

import (
	"github.com/tyler-smith/go-bip39/wordlists"
)

// Wordlists provides access to BIP39 wordlists in various languages.
// These are used for generating and validating mnemonic phrases.
var Wordlists = struct {
	// English is the English BIP39 wordlist (2048 words).
	English []string
	// Japanese is the Japanese BIP39 wordlist.
	Japanese []string
	// ChineseSimplified is the Simplified Chinese BIP39 wordlist.
	ChineseSimplified []string
	// ChineseTraditional is the Traditional Chinese BIP39 wordlist.
	ChineseTraditional []string
	// Czech is the Czech BIP39 wordlist.
	Czech []string
	// French is the French BIP39 wordlist.
	French []string
	// Italian is the Italian BIP39 wordlist.
	Italian []string
	// Korean is the Korean BIP39 wordlist.
	Korean []string
	// Spanish is the Spanish BIP39 wordlist.
	Spanish []string
	// Portuguese is the Portuguese BIP39 wordlist.
	Portuguese []string
}{
	English:            wordlists.English,
	Japanese:           wordlists.Japanese,
	ChineseSimplified:  wordlists.ChineseSimplified,
	ChineseTraditional: wordlists.ChineseTraditional,
	Czech:              wordlists.Czech,
	French:             wordlists.French,
	Italian:            wordlists.Italian,
	Korean:             wordlists.Korean,
	Spanish:            wordlists.Spanish,
	// Note: Portuguese is not available in the standard go-bip39 package
	// Using English as a fallback
	Portuguese: wordlists.English,
}

// GetWordlist returns a wordlist by language name.
//
// Supported languages: "english", "japanese", "chinese_simplified",
// "chinese_traditional", "czech", "french", "italian", "korean", "spanish"
func GetWordlist(language string) ([]string, error) {
	switch language {
	case "english", "en":
		return Wordlists.English, nil
	case "japanese", "ja":
		return Wordlists.Japanese, nil
	case "chinese_simplified", "zh-hans", "zh_hans":
		return Wordlists.ChineseSimplified, nil
	case "chinese_traditional", "zh-hant", "zh_hant":
		return Wordlists.ChineseTraditional, nil
	case "czech", "cs":
		return Wordlists.Czech, nil
	case "french", "fr":
		return Wordlists.French, nil
	case "italian", "it":
		return Wordlists.Italian, nil
	case "korean", "ko":
		return Wordlists.Korean, nil
	case "spanish", "es":
		return Wordlists.Spanish, nil
	default:
		return nil, ErrInvalidWordlist
	}
}
