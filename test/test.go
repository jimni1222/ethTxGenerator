package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"os/exec"
	"strings"
)

var (
	exeFileName string
)

// You can test how to consume what this project made with `go run ./test/test.go`
func main2() {
	exeFileName = "ethTxGenerator"
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path = path + "/" + exeFileName
	_, err = exec.LookPath(path)
	if err != nil {
		fmt.Println("executable file is not existed.")
	}

	endpoint := "http://3.36.94.96:8551"
	// endpoint := "https://api.baobab.klaytn.net:8651"
	txType := "1"
	chainID := "10000"
	// chainID := "1001"
	// gasPrice := "25000000000"
	gasPrice := "0x5d21dba00"
	gas := "200000"
	// gas := "0xdbba0"
	baseFee := "0"
	value := "0"
	// value := "1"
	from := "0xd436ae3aeb4498ae1359b1176f4ce0431a360870964248887556489934a05377"
	// from := "0x0c66be6f6a0c539be9e99d883e897cdc10cb016a958c9a61485ae57cace5d7bf"
	nonce := "0"
	// nonce := "919"
	// to := "0xabc4da09ab28ade7fe02e34f7a94f3a0bcafbc05"
	// to := "random"
	to := "nil"
	// code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
	// input := "0x6353586b00000000000000000000000042f2f89f7d8061dd99d41b46ada78bb9ccaf4e77"
	input := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
	//
	// ctx := context.Background()
	// rpcClient, err := rpc.DialContext(ctx, endpoint)
	// if err != nil {
	//	fmt.Print("Failed to connect RPC client: %v", err)
	//	os.Exit(1)
	// }
	// geth := gethclient.New(rpcClient)
	// addr := common.HexToAddress("0xabc4da09ab28ade7fe02e34f7a94f3a0bcafbc05")
	// msgCall := ethereum.CallMsg{
	//	From:  addr,
	//	To:    &addr,
	//	Value: new(big.Int).SetUint64(1),
	// }
	// accessList, _, _, err := geth.CreateAccessList(ctx, msgCall)
	//
	// if err != nil {
	//	fmt.Printf("failed to create access list: %v", err)
	// }
	// fmt.Print(accessList)

	cmd := exec.Command(path, endpoint, txType, chainID, gasPrice, gas, baseFee, value, from, nonce, to, input)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", stdoutStderr)

	// s := "0xf4c110e7d8e66a566d4b70cbec4cb0871f57623550a6b5222b07b97614f86bf5\n"
	// s = strings.TrimSuffix(s, "\n")
	// fmt.Println(s)
	// if len(s) > 1 {
	//	if s[0:2] == "0x" || s[0:2] == "0X" {
	//		s = s[2:]
	//	}
	// }
	// if len(s)%2 == 1 {
	//	s = "0" + s
	// }
	//
	// fmt.Println(s)
	//
	// h, _ := hex.DecodeString(s)
	// fmt.Println(h)
	//
	// hash := common.BytesToHash(h)
	// fmt.Println(hash.String())
}

// You can test how to consume what this project made with `go run ./test/test.go`
func main() {
	exeFileName = "ethTxGenerator"
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path = path + "/" + exeFileName
	_, err = exec.LookPath(path)
	if err != nil {
		fmt.Println("executable file is not existed.")
	}
	fmt.Printf("Path: %v\n", path)

	endpoint := "http://52.78.55.68:8551"
	// endpoint := "https://api.baobab.klaytn.net:8651"
	gasPrice := "0x5d21dba00"
	gas := "200000"
	value := "0"
	from := "0x926601ab01b1a70ec84af781aeca5e71575e08e8"
	to := "0xd109065ee17e2dc20b3472a4d4fb5907bd687d09"
	data, err := GenerateData(from)
	if err != nil {
		fmt.Println("failed to encode data")
	}
	input := hex.EncodeToString(data)

	cmd := exec.Command(path, endpoint, "eth_call", from, to, gasPrice, gas, value, input)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func GenerateData(addrString string) ([]byte, error) {
	addr := common.HexToAddress(addrString)
	abii, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	if err != nil {
		return nil, fmt.Errorf("failed to abi.JSON: %v", err)
	}
	data, err := abii.Pack("balanceOf", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to abi.Pack: %v", err)
	}
	return data, nil
}
