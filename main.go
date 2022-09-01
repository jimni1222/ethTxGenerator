package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var (
	url    string
	client *ethclient.Client
	geth   *gethclient.Client

	txType   int
	chainID  *big.Int
	gasPrice *big.Int
	gas      *big.Int
	baseFee  *big.Int
	value    *big.Int
	data     []byte

	fromAccount *Account
	nonce       uint64
	fromAddress *common.Address
	to          *common.Address
	blockNumber *big.Int

	ctx context.Context

	invalidArgsForSendingTx  = "invalid arguments: endpoint txType chainID gasPrice gas baseFee value fromPrivateKey nonce toAddress [data]. "
	invalidArgsForCallingMsg = "invalid arguments: endpoint eth_call fromAddress toAddress gasPrice gas value data [, blockNumber]."
)

type Account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

// `go build` will generate executable file named "ethTxGenerator"
func main() {
	if len(os.Args) < 2 {
		fmt.Print("invalid arguments: no arguments are passed.")
		os.Exit(1)
	}
	// os.Args[0] will be program path
	args := os.Args[1:]

	url = args[0]
	// fmt.Print("Test url: ", url)
	ctx = context.Background()
	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		fmt.Print("Failed to connect RPC client: %v", err)
		os.Exit(1)
	}
	client = ethclient.NewClient(rpcClient)
	geth = gethclient.New(rpcClient)

	if args[1] == "eth_call" {
		// In case of "eth_call", since eth_call comes in as a program argument,
		// it is excluded and transmitted as a function parameter.
		callEthMsg(args[2:])
		return
	}

	// If the method name is not entered as a parameter,
	// send Transaction is operated by default.
	sendEthTx(args)
}

func callEthMsg(args []string) {
	argsLen := len(args)

	if argsLen < 6 {
		fmt.Print(invalidArgsForCallingMsg, "not enough arguments to call an ethereum message.")
		os.Exit(1)
	}

	fromAddress = parseToAddress(args[0])
	// fmt.Print("Test fromAddress: ", fromAddress.String())

	to = parseToAddress(args[1])
	// fmt.Print("Test to account: ", to.String()

	gas = parseToBigInt(args[2], "gas")
	// fmt.Print("Test gas: ", gas)

	gasPrice = parseToBigInt(args[3], "gasPrice")
	// fmt.Print("Test gas price: ", gasPrice)

	value = parseToBigInt(args[4], "value")
	// fmt.Print("Test value: ", value)

	data = common.FromHex(args[5])
	// fmt.Print("Test data: ", data)

	// `blockNumber` is optional field, so if user pass the last parameter, then set to `blockNumber`.
	if argsLen == 8 {
		blockNumber = parseToBigInt(args[6], "blockNumber")
		// fmt.Print("Test blockNumber: ", blockNumber)
	} else {
		blockNumber = nil
	}

	callMsg := ethereum.CallMsg{
		From:     *fromAddress,
		To:       to,
		Gas:      gas.Uint64(),
		GasPrice: gasPrice,
		Data:     data,
		Value:    value,
	}
	// fmt.Printf("From: %v / To: %v / Gas: %v / GasPrice: %v / Data: %v / Value: %v \n", callMsg.From.String(), callMsg.To.String(), gas, gasPrice, hex.EncodeToString(data), value.String())

	var ret []byte
	ret, err := client.CallContract(ctx, callMsg, blockNumber)
	if err != nil {
		fmt.Print("Error: " + err.Error())
		os.Exit(1)
	}

	resultString := hex.EncodeToString(ret)
	fmt.Print(resultString)
}

func sendEthTx(args []string) {
	argsLen := len(args)

	if argsLen < 10 {
		fmt.Print(invalidArgsForSendingTx, "not enough arguments to send an ethereum transaction.")
		os.Exit(1)
	}

	url = args[0]
	// fmt.Print("Test url: ", url)
	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		fmt.Print("Failed to connect RPC client: %v", err)
		os.Exit(1)
	}
	client = ethclient.NewClient(rpcClient)
	geth = gethclient.New(rpcClient)

	txType, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Print(invalidArgsForSendingTx, err)
		os.Exit(1)
	}
	// fmt.Print("Test tx type: ", txType)

	chainID = parseToBigInt(args[2], "chainID")
	// fmt.Print("Test chain id: ", chainID)

	gasPrice = parseToBigInt(args[3], "gasPrice")
	// fmt.Print("Test gas price: ", gasPrice)

	gas = parseToBigInt(args[4], "gas")
	// fmt.Print("Test gas: ", gas)

	baseFee = parseToBigInt(args[5], "baseFee")
	// fmt.Print("Test base fee: ", baseFee)

	value = parseToBigInt(args[6], "value")
	// fmt.Print("Test value: ", value)

	fromAccount = createTestAccountWithPrivateKey(args[7])
	// fmt.Print("Test from account: ", from.address)

	nonce, err = strconv.ParseUint(args[8], 10, 0)
	if err != nil {
		fmt.Print(invalidArgsForSendingTx, err)
		os.Exit(1)
	}
	// fmt.Print("Test nonce: ", nonce)

	to = parseToAddress(args[9])
	// fmt.Print("Test to account: ", to.String())

	// `data` is optional field, so if user pass the last parameter, then set to `data`.
	if argsLen == 11 {
		data = common.FromHex(args[10])
	} else {
		data = []byte{}
	}

	tx := createTxWithGeth()
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	hash := tx.Hash().String()
	fmt.Print(hash)
}

func parseToBigInt(arg string, paramName string) *big.Int {
	i := new(big.Int)
	base := 10
	if strings.Contains(arg, "0x") {
		arg = arg[2:]
		base = 16
	}
	i, ok := i.SetString(arg, base)
	if !ok {
		fmt.Print("Invalid parameter: ", paramName)
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
		fmt.Print("invalid private key: ", prv)
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

	// Create messsage call
	maxPriorityFeePerGas := gasPrice
	maxFeePerGas := big.NewInt(0).Add(big.NewInt(0).Mul(baseFee, big.NewInt(2)), maxPriorityFeePerGas)
	msgCall := ethereum.CallMsg{
		From:      fromAccount.address,
		To:        to,
		Gas:       gas.Uint64(),
		GasPrice:  gasPrice,
		Data:      data,
		Value:     value,
		GasFeeCap: maxFeePerGas,
		GasTipCap: maxPriorityFeePerGas,
	}

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
		accessList, _, _, err := geth.CreateAccessList(context.Background(), msgCall)
		if err != nil {
			fmt.Printf("failed to create access list: %v", err)
			os.Exit(1)
		}
		txdata = &types.AccessListTx{
			ChainID:    chainID,
			Nonce:      nonce,
			To:         to,
			Gas:        gas.Uint64(),
			GasPrice:   gasPrice,
			AccessList: *accessList,
			Data:       data,
			Value:      value,
		}
	} else if txType == 2 {
		accessList, _, _, err := geth.CreateAccessList(context.Background(), msgCall)
		if err != nil {
			fmt.Printf("failed to create access list: %v", err)
			os.Exit(1)
		}
		txdata = &types.DynamicFeeTx{
			ChainID:    chainID,
			Nonce:      nonce,
			To:         to,
			GasFeeCap:  maxFeePerGas,
			GasTipCap:  maxPriorityFeePerGas,
			Gas:        gas.Uint64(),
			AccessList: *accessList,
			Data:       data,
			Value:      value,
		}
	} else {
		fmt.Printf("invalid tx type: %v", txType)
		os.Exit(1)
	}

	tx, err := types.SignNewTx(fromAccount.privateKey, signer, txdata)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return tx
}
