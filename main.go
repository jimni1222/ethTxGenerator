package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var (
	txType   int
	chainID  *big.Int
	gasPrice *big.Int
	gas      *big.Int
	baseFee  *big.Int

	from  *Account
	nonce int
	to    common.Address

	invalidArgs string
)

type Account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

// `go build` will generate executable file named "ethTxGenerator"
func main() {
	invalidArgs = "invalid arguments: txType, chainID, gasPrice, gas, baseFee, from private key, nonce and to address should be passed as arguments. "

	if len(os.Args) < 2 {
		fmt.Println(invalidArgs, "no arguments are passed.")
		os.Exit(1)
	}
	// os.Args[0] will be program path
	args := os.Args[1:]

	if len(args) < 8 {
		fmt.Println(invalidArgs, "not enough arguments.")
		os.Exit(1)
	}

	txType, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(invalidArgs, err)
		os.Exit(1)
	}
	fmt.Println("Test tx type: ", txType)

	chainID = parseToBigInt(args[1])
	fmt.Println("Test chain id: ", chainID)

	gasPrice = parseToBigInt(args[2])
	fmt.Println("Test gas price: ", gasPrice)

	gas = parseToBigInt(args[3])
	fmt.Println("Test gas: ", gas)

	baseFee = parseToBigInt(args[4])
	fmt.Println("Test base fee: ", baseFee)

	from = createTestAccountWithPrivateKey(args[5])
	fmt.Println("Test from account: ", from.address)

	nonce, err = strconv.Atoi(args[6])
	if err != nil {
		fmt.Println(invalidArgs, err)
		os.Exit(1)
	}
	fmt.Println("Test nonce: ", nonce)

	to = parseToAddress(args[7])
	fmt.Println("Test to account: ", to.String())

	createTxWithGeth()
}

func parseToBigInt(arg string) *big.Int {
	i := new(big.Int)
	base := 10
	if strings.Contains(arg, "0x") {
		arg = arg[2:]
		base = 16
	}
	i, ok := i.SetString(arg, base)
	if !ok {
		fmt.Println(invalidArgs)
		os.Exit(1)
	}
	return i
}

func parseToAddress(addr string) common.Address {
	if addr == "random" {
		return generateRandomAccount().address
	}

	return common.HexToAddress(addr)
}

func createTestAccountWithPrivateKey(prv string) *Account {
	if prv == "random" {
		return generateRandomAccount()
	}

	acc, err := crypto.HexToECDSA(prv)
	if err != nil {
		fmt.Println("invalid private key: ", prv)
		os.Exit(1)
	}

	return createAccount(acc)
}

func createAccount(k *ecdsa.PrivateKey) *Account {

	return &Account{
		address:    crypto.PubkeyToAddress(k.PublicKey),
		privateKey: k,
	}
}

func generateRandomAccount() *Account {
	return createAccount(genRandomPrivateKey())
}

func genRandomPrivateKey() *ecdsa.PrivateKey {
	acc, err := crypto.GenerateKey()
	if err != nil {
		os.Exit(1)
	}

	return acc
}

func createTxWithGeth() string {
	var txdata types.TxData
	signer := types.NewLondonSigner(chainID)
	nonce := uint64(0)

	if txType == 0 {
		txdata = &types.LegacyTx{
			Nonce:    nonce,
			To:       &to,
			Gas:      30000,
			GasPrice: gasPrice,
			Data:     []byte{},
		}
	} else if txType == 1 {
		txdata = &types.AccessListTx{
			ChainID:  chainID,
			Nonce:    nonce,
			To:       &to,
			Gas:      30000,
			GasPrice: gasPrice,
			AccessList: types.AccessList{
				{
					Address:     to,
					StorageKeys: []common.Hash{{0}},
				},
			},
			Data: []byte{},
		}
	} else if txType == 2 {
		maxPriorityFeePerGas := gasPrice
		maxFeePerGas := big.NewInt(0).Add(big.NewInt(0).Mul(baseFee, big.NewInt(2)), maxPriorityFeePerGas)

		txdata = &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &to,
			GasFeeCap: maxFeePerGas,
			GasTipCap: maxPriorityFeePerGas,
			Gas:       30000,
			AccessList: types.AccessList{
				{
					Address:     to,
					StorageKeys: []common.Hash{{0}},
				},
			},
			Data: []byte{},
		}
	} else {
		fmt.Println("invalid tx type: ", txType)
	}

	tx, err := types.SignNewTx(from.privateKey, signer, txdata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encoded, err := tx.MarshalBinary()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Signed Tx")
	fmt.Println(hex.EncodeToString(encoded))

	return hex.EncodeToString(encoded)

}
