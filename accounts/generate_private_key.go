package accounts

import (
	"github.com/ChefBingbong/viem-go/crypto"
	"github.com/ChefBingbong/viem-go/utils"
)

func GeneratePrivateKey() string {
	privateKey, err := crypto.NewPrivKey().MarshalJSON()
	if err != nil {
		panic(err)
	}
	return utils.BytesToHex(privateKey)
}
