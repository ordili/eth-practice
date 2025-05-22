package tx

import (
	"eth-practice/util"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)

func TestGetBlock(t *testing.T) {
	client := util.GetClient()
	GetBlock(client, nil)
}

func TestTransferETH(t *testing.T) {
	key := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privateKey, _ := crypto.HexToECDSA(key)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	value := big.NewInt(1000000000000000000)
	client := util.GetClient()

	TransferETH(client, privateKey, toAddress, value)
}

func TestCreateRawTransaction(t *testing.T) {
	client := util.GetClient()
	key := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privateKey, _ := crypto.HexToECDSA(key)
	to := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	value := big.NewInt(1000000000000000000)
	rawTx := CreateRawTransaction(client, privateKey, to, value)
	fmt.Println("rawTx is  ", rawTx)
	SendRawTransaction(client, rawTx)
}
