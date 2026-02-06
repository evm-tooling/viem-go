package accounts_test

import (
	"math/big"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/accounts"
	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// Test private key (Anvil account 0)
const testPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const testAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

// Test mnemonic
const testMnemonic = "test test test test test test test test test test test junk"

var _ = Describe("Accounts", func() {
	Describe("GeneratePrivateKey", func() {
		It("should generate a valid private key", func() {
			privateKey := accounts.GeneratePrivateKey()
			Expect(privateKey).To(HavePrefix("0x"))
			Expect(len(privateKey)).To(Equal(66)) // 0x + 64 hex chars
		})

		It("should generate unique keys each time", func() {
			key1 := accounts.GeneratePrivateKey()
			key2 := accounts.GeneratePrivateKey()
			Expect(key1).NotTo(Equal(key2))
		})
	})

	Describe("GenerateMnemonic", func() {
		It("should generate a 12-word mnemonic by default", func() {
			mnemonic, err := accounts.GenerateMnemonic()
			Expect(err).NotTo(HaveOccurred())
			words := strings.Fields(mnemonic)
			Expect(len(words)).To(Equal(12))
		})

		It("should generate a 24-word mnemonic with Mnemonic256", func() {
			mnemonic, err := accounts.GenerateMnemonic(accounts.GenerateMnemonicOptions{
				Strength: accounts.Mnemonic256,
			})
			Expect(err).NotTo(HaveOccurred())
			words := strings.Fields(mnemonic)
			Expect(len(words)).To(Equal(24))
		})

		It("should generate valid mnemonics", func() {
			mnemonic, err := accounts.GenerateMnemonic()
			Expect(err).NotTo(HaveOccurred())
			Expect(accounts.ValidateMnemonic(mnemonic)).To(BeTrue())
		})
	})

	Describe("ValidateMnemonic", func() {
		It("should validate a valid mnemonic", func() {
			Expect(accounts.ValidateMnemonic(testMnemonic)).To(BeTrue())
		})

		It("should reject an invalid mnemonic", func() {
			Expect(accounts.ValidateMnemonic("invalid mnemonic phrase")).To(BeFalse())
		})
	})

	Describe("PrivateKeyToAccount", func() {
		It("should create an account from a private key", func() {
			account, err := accounts.PrivateKeyToAccount(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(account.GetAddress()).To(BeEquivalentTo(testAddress))
			Expect(account.GetType()).To(Equal(accounts.AccountTypeLocal))
			Expect(account.GetSource()).To(Equal(accounts.AccountSourcePrivateKey))
		})

		It("should return public key", func() {
			account, err := accounts.PrivateKeyToAccount(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(account.GetPublicKey()).To(HavePrefix("0x04"))
		})

		It("should fail with invalid private key", func() {
			_, err := accounts.PrivateKeyToAccount("0xinvalid")
			Expect(err).To(HaveOccurred())
		})

		It("should sign messages", func() {
			account, err := accounts.PrivateKeyToAccount(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())

			sig, err := account.SignMessage(signature.NewSignableMessage("hello world"))
			Expect(err).NotTo(HaveOccurred())
			Expect(sig).To(HavePrefix("0x"))
			Expect(len(sig)).To(Equal(132))
		})

		It("should sign transactions", func() {
			account, err := accounts.PrivateKeyToAccount(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())

			tx := &transaction.Transaction{
				Type:                 transaction.TransactionTypeEIP1559,
				ChainId:              1,
				Nonce:                0,
				MaxPriorityFeePerGas: big.NewInt(1000000000),
				MaxFeePerGas:         big.NewInt(2000000000),
				Gas:                  big.NewInt(21000),
				To:                   "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
				Value:                big.NewInt(1000000000000000000),
			}

			signedTx, err := account.SignTransaction(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(signedTx).To(HavePrefix("0x02"))
		})

		It("should sign typed data", func() {
			account, err := accounts.PrivateKeyToAccount(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())

			typedData := signature.TypedDataDefinition{
				Domain: signature.TypedDataDomain{
					Name:              "Test",
					Version:           "1",
					ChainId:           big.NewInt(1),
					VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
				},
				Types: map[string][]signature.TypedDataField{
					"Message": {
						{Name: "content", Type: "string"},
					},
				},
				PrimaryType: "Message",
				Message: map[string]any{
					"content": "Hello",
				},
			}

			sig, err := account.SignTypedData(typedData)
			Expect(err).NotTo(HaveOccurred())
			Expect(sig).To(HavePrefix("0x"))
		})

		It("should sign authorizations", func() {
			account, err := accounts.PrivateKeyToAccount(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())

			signedAuth, err := account.SignAuthorization(accounts.AuthorizationRequest{
				Address: "0x1234567890123456789012345678901234567890",
				ChainId: 1,
				Nonce:   0,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(signedAuth.Address).To(Equal("0x1234567890123456789012345678901234567890"))
			Expect(signedAuth.R).To(HavePrefix("0x"))
			Expect(signedAuth.S).To(HavePrefix("0x"))
		})
	})

	Describe("MnemonicToAccount", func() {
		It("should create an account from a mnemonic", func() {
			account, err := accounts.MnemonicToAccount(testMnemonic)
			Expect(err).NotTo(HaveOccurred())
			Expect(account.GetType()).To(Equal(accounts.AccountTypeLocal))
		})

		It("should fail with invalid mnemonic", func() {
			_, err := accounts.MnemonicToAccount("invalid mnemonic")
			Expect(err).To(HaveOccurred())
		})

		It("should derive different addresses with different account indices", func() {
			account0, err := accounts.MnemonicToAccount(testMnemonic, accounts.MnemonicToAccountOptions{
				HDOptions: accounts.HDOptions{AccountIndex: 0},
			})
			Expect(err).NotTo(HaveOccurred())

			account1, err := accounts.MnemonicToAccount(testMnemonic, accounts.MnemonicToAccountOptions{
				HDOptions: accounts.HDOptions{AccountIndex: 1},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(account0.GetAddress()).NotTo(Equal(account1.GetAddress()))
		})

		It("should support passphrase", func() {
			accountWithoutPassphrase, _ := accounts.MnemonicToAccount(testMnemonic)
			accountWithPassphrase, _ := accounts.MnemonicToAccount(testMnemonic, accounts.MnemonicToAccountOptions{
				Passphrase: "my-secret-passphrase",
			})

			Expect(accountWithoutPassphrase.GetAddress()).NotTo(Equal(accountWithPassphrase.GetAddress()))
		})
	})

	Describe("ToAccount", func() {
		It("should create a json-rpc account from address string", func() {
			account, err := accounts.ToAccountFromAddress(testAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(account.Address).To(Equal(testAddress))
			Expect(account.Type).To(Equal(accounts.AccountTypeJSONRPC))
		})

		It("should create a local account from custom source", func() {
			account, err := accounts.ToAccount(accounts.CustomSource{
				Address: testAddress,
				SignMessage: func(message signature.SignableMessage) (string, error) {
					return "0x1234", nil
				},
				SignTransaction: func(tx *transaction.Transaction) (string, error) {
					return "0x5678", nil
				},
				SignTypedData: func(data signature.TypedDataDefinition) (string, error) {
					return "0x9abc", nil
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(account.Addr).To(Equal(testAddress))
			Expect(account.Type).To(Equal(accounts.AccountTypeLocal))
			Expect(account.Source).To(Equal(accounts.AccountSourceCustom))
		})

		It("should fail with invalid address", func() {
			_, err := accounts.ToAccountFromAddress("0xinvalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Wordlists", func() {
		It("should provide English wordlist", func() {
			Expect(len(accounts.Wordlists.English)).To(Equal(2048))
		})

		It("should get wordlist by name", func() {
			wordlist, err := accounts.GetWordlist("english")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(wordlist)).To(Equal(2048))
		})

		It("should fail for unknown wordlist", func() {
			_, err := accounts.GetWordlist("unknown")
			Expect(err).To(HaveOccurred())
		})
	})
})
