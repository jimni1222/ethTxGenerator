package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var (
	url string

	txType   int
	chainID  *big.Int
	gasPrice *big.Int
	gas      *big.Int
	baseFee  *big.Int
	value    *big.Int
	data     []byte

	from  *Account
	nonce uint64
	to    *common.Address

	invalidArgs string
)

type Account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

// `go build` will generate executable file named "ethTxGenerator"
func main() {
	invalidArgs = "invalid arguments: endpoint txType chainID gasPrice gas baseFee value fromPrivateKey nonce toAddress [data]. "

	if len(os.Args) < 2 {
		fmt.Println(invalidArgs, "no arguments are passed.")
		os.Exit(1)
	}
	// os.Args[0] will be program path
	args := os.Args[1:]
	argsLen := len(args)

	if argsLen < 10 {
		fmt.Println(invalidArgs, "not enough arguments.")
		os.Exit(1)
	}

	url = args[0]
	//fmt.Println("Test url: ", url)
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Failed to connect Eth RPC client: %v", err)
		os.Exit(1)
	}

	txType, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(invalidArgs, err)
		os.Exit(1)
	}
	//fmt.Println("Test tx type: ", txType)

	chainID = parseToBigInt(args[2])
	//fmt.Println("Test chain id: ", chainID)

	gasPrice = parseToBigInt(args[3])
	//fmt.Println("Test gas price: ", gasPrice)

	gas = parseToBigInt(args[4])
	//fmt.Println("Test gas: ", gas)

	baseFee = parseToBigInt(args[5])
	//fmt.Println("Test base fee: ", baseFee)

	value = parseToBigInt(args[6])
	//fmt.Println("Test value: ", value)

	from = createTestAccountWithPrivateKey(args[7])
	//fmt.Println("Test from account: ", from.address)

	nonce, err = strconv.ParseUint(args[8], 10, 0)
	if err != nil {
		fmt.Println(invalidArgs, err)
		os.Exit(1)
	}
	//fmt.Println("Test nonce: ", nonce)

	to = parseToAddress(args[9])
	//fmt.Println("Test to account: ", to.String())

	// `data` is optional field, so if user pass the last parameter, then set to `data`.
	if argsLen == 11 {
		data = common.FromHex(args[10])
	} else {
		data = []byte{}
	}

	tx := createTxWithGeth()
	ctx := context.Background()
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hash := tx.Hash().String()
	fmt.Println(hash)
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

func parseToAddress(addr string) *common.Address {
	if addr == "random" {
		return &generateRandomAccount().address
	}
	if addr == "nil" {
		return nil
	}

	address := common.HexToAddress(addr)
	return &address
}

func createTestAccountWithPrivateKey(prv string) *Account {
	if prv == "random" {
		return generateRandomAccount()
	}

	if strings.Contains(prv, "0x") {
		prv = prv[2:]
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

func createTxWithGeth() *types.Transaction {
	var txdata types.TxData
	signer := types.NewLondonSigner(chainID)

	if txType == 0 {
		txdata = &types.LegacyTx{
			Nonce:    nonce,
			To:       to,
			Gas:      gas.Uint64(),
			GasPrice: gasPrice,
			Data:     data,
			Value:    value,
		}
	} else if txType == 1 {
		txdata = &types.AccessListTx{
			ChainID:  chainID,
			Nonce:    nonce,
			To:       to,
			Gas:      gas.Uint64(),
			GasPrice: gasPrice,
			AccessList: types.AccessList{
				{
					Address:     *to,
					StorageKeys: []common.Hash{{0}},
				},
			},
			Data:  data,
			Value: value,
		}
	} else if txType == 2 {
		maxPriorityFeePerGas := gasPrice
		maxFeePerGas := big.NewInt(0).Add(big.NewInt(0).Mul(baseFee, big.NewInt(2)), maxPriorityFeePerGas)

		txdata = &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        to,
			GasFeeCap: maxFeePerGas,
			GasTipCap: maxPriorityFeePerGas,
			Gas:       gas.Uint64(),
			AccessList: types.AccessList{
				{
					Address:     *to,
					StorageKeys: []common.Hash{{0}},
				},
			},
			Data:  data,
			Value: value,
		}
	} else {
		fmt.Println("invalid tx type: ", txType)
	}

	tx, err := types.SignNewTx(from.privateKey, signer, txdata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return tx
}
