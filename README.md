# ethTxGenerator
Generate an ethereum transaction via Geth.

## Make an executable file
```shell
go build
```


## Execute ethTxGenerator
You need to pass some arguments needed to make a transaction.

```shell
./ethTxGenerator txType chainID gasPrice gas baseFee fromPrivateKey nonce to
```

For example, you can execute like below.

```shell
./ethTxGenerator 1 8217 0x5d21dba00 0xdbba0 0 0xFromPrivateKey 0 0xToAddress
```

You can send either int or hex string for `chainId`, `gasPrice`, `gas` and `baseFee`.


If you want to make a tx with a random account, you can pass `"random"` instead of `from` and `to` like below.

```shell
./ethTxGenerator txType chainID gasPrice gas baseFee random nonce random
```