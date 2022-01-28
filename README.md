# ethTxGenerator
Generate an ethereum transaction via Geth.

## Make an executable file
```shell
go build
```


## Execute ethTxGenerator
You need to pass some arguments needed to make a transaction.

```shell
./ethTxGenerator endPoint txType chainID gasPrice gas baseFee value fromPrivateKey nonce to [data]
```

For example, you can execute like below without data.

```shell
./ethTxGenerator http://your_en_url:port 1 8217 0x5d21dba00 0xdbba0 0 1 0xFromPrivateKey 0 0xToAddress
```

You can send either int or hex string for `chainId`, `gasPrice`, `gas`, `baseFee` and `value`.

If you want to deploy a smart contract, then you can pass `nil` string to `to` argument like below.
Note that you need to pass data argument together.

```shell
./ethTxGenerator http://your_en_url:port 1 8217 0x5d21dba00 0xdbba0 0 1 0xFromPrivateKey 0 nil 0xByteCode
```

If you want to execute a smart contract, then you can pass a smart contract address to `to` argument and an encoded function call to `data` argument like below.

```shell
./ethTxGenerator http://your_en_url:port 1 8217 0x5d21dba00 0xdbba0 0 1 0xFromPrivateKey 0 0xSmartContractAddress 0xEncodedFunctionCall
```

## For easy testing

If you want to make a tx with a random account, you can pass `"random"` instead of `from` and `to` like below.
Maybe this can be used only to check a transaction format which is generated from this ethTxGenerator.

```shell
./ethTxGenerator endPoint txType chainID gasPrice gas baseFee value random nonce random
```