package tx

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"eth-practice/util"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"math/big"
)

func GetBlockHeader(client *ethclient.Client) *types.Header {

	number := big.NewInt(8378757)
	header, err := client.HeaderByNumber(context.Background(), number)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String())
	return header
}

func GetBlock(client *ethclient.Client, blockNumber *big.Int) *types.Block {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())                // 5671744
	fmt.Println(block.Time())                           // 1527211625
	fmt.Println(block.Difficulty().Uint64())            // 3217000136609065
	fmt.Println(block.Hash().Hex())                     // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println("tx len {}", len(block.Transactions())) // 144

	for _, tx := range block.Transactions()[:1] {
		fmt.Println(tx.Hash().Hex())
		fmt.Println(tx.Value().String())    // 10000000000000000
		fmt.Println(tx.Gas())               // 105000
		fmt.Println(tx.GasPrice().Uint64()) // 102000000000
		fmt.Println(tx.Nonce())             // 110644
		fmt.Println(tx.Data())              // []
		fmt.Println(tx.To().Hex())

		//if msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId())); err == nil {
		//	fmt.Println(msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		//}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status)
	}
	return block
}

func TransferETH(client *ethclient.Client, privateKey *ecdsa.PrivateKey, toAddress common.Address, value *big.Int) {

	fromAddress, _ := util.GetPublicAddressFromPrivateAddress(privateKey)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s \n", signedTx.Hash().Hex())
}

func TransferToken(client *ethclient.Client) {

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256(transferFnSignature)
	methodID := []byte("0x" + hex.EncodeToString(hash[:4]))

	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit) // 23256

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit*2, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
}

func BlockSubscribe() {
	client, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64())   // 3477413
			fmt.Println(block.Time())              // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		}
	}
}

func CreateRawTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, toAddress common.Address, value *big.Int) string {

	fromAddress, _ := util.GetPublicAddressFromPrivateAddress(privateKey)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	rlpBytes, err := signedTx.MarshalBinary()

	if err != nil {
		fmt.Printf("Failed to encode transaction: %v\n", err)
		return ""
	}

	// Convert to hex string
	rlpHex := hex.EncodeToString(rlpBytes)
	fmt.Println("RLP Encoding:", rlpHex)

	txHash := tx.Hash().Hex()
	fmt.Println("Transaction Hash:", txHash)

	return rlpHex
}

func SendRawTransaction(client *ethclient.Client, rawTx string) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		log.Fatal(err)
	}
	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
