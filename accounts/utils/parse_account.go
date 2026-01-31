package utils

// ParseAccount parses an account from either an address string or an Account struct.
// If given a string address, it returns a JSON-RPC account.
// If given an Account, it returns it as-is.
//
// Example:
//
//	// From address string
//	account := ParseAccount("0x1234567890123456789012345678901234567890")
//	// Returns: Account{Address: "0x1234...", Type: AccountTypeJSONRPC}
//
//	// From Account struct
//	account := ParseAccountFromAccount(Account{Address: "0x1234...", Type: AccountTypeLocal})
//	// Returns: Account{Address: "0x1234...", Type: AccountTypeLocal}
func ParseAccount(address string) Account {
	return Account{
		Address: address,
		Type:    AccountTypeJSONRPC,
	}
}

// ParseAccountFromAccount returns the account as-is.
// This is provided for API consistency with the TypeScript version.
func ParseAccountFromAccount(account Account) Account {
	return account
}

// ParseAccountGeneric parses an account from either a string or Account.
// This mirrors the TypeScript overloaded function behavior.
func ParseAccountGeneric(accountOrAddress any) Account {
	switch v := accountOrAddress.(type) {
	case string:
		return ParseAccount(v)
	case Account:
		return v
	case *Account:
		if v != nil {
			return *v
		}
		return Account{}
	default:
		return Account{}
	}
}
