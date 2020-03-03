package transactions

import (
	"encoding/base64"
	"fmt"

	sdkTxs "github.com/SebastianJ/harmony-sdk/transactions"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/numeric"
)

// IsTransactionSuccessful - checks if a transaction is successful given a transaction response
func IsTransactionSuccessful(txResponse map[string]interface{}) (success bool) {
	txStatus, ok := txResponse["status"].(string)

	if txStatus != "" && ok {
		success = (txStatus == "0x1")
	}

	return success
}

// SendTransaction - send transactions
func SendTransaction(keystore *keystore.KeyStore, account *accounts.Account, rpcClient *rpc.HTTPMessenger, chain *common.ChainID, fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount numeric.Dec, gasLimit int64, gasPrice numeric.Dec, currentNonce uint64, txData string, passphrase string, node string, timeout int) (map[string]interface{}, error) {
	fmt.Println(fmt.Sprintf("Sending tx - From: %s, From Shard: %d, To: %s, To Shard: %d, Amount: %s, Nonce: %d, Node: %s, RPCClient: %+v", fromAddress, fromShardID, toAddress, toShardID, amount, currentNonce, node, rpcClient))

	if txData != "" {
		txData = base64.StdEncoding.EncodeToString([]byte(txData))
	}

	txResult, err := sdkTxs.SendTransaction(keystore, account, rpcClient, chain, fromAddress, fromShardID, toAddress, toShardID, amount, gasLimit, gasPrice, currentNonce, txData, passphrase, node, timeout)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error occurred: %s", err.Error()))
		return nil, err
	}

	return txResult, nil
}
