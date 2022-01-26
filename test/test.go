package main

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	exeFileName string
)

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

	endpoint := "url"
	txType := "0"
	chainID := "1001"
	gasPrice := "0x5d21dba00"
	gas := "0xdbba0"
	baseFee := "0"
	value := "1"
	from := "0x{private key}"
	nonce := "0"
	to := "random"

	cmd := exec.Command(path, endpoint, txType, chainID, gasPrice, gas, baseFee, value, from, nonce, to)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
