package transactions

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	"github.com/SebastianJ/harmony-sdk/transactions"
	"github.com/SebastianJ/harmony-stress/utils"
	"github.com/harmony-one/harmony/numeric"
)

// BulkSendTransactions - sends transactions in bulk synchronously
func BulkSendTransactions() (txResults []map[string]interface{}, err error) {
	currentNonce := Configuration.Account.Nonce
	gasPrice := Configuration.Transactions.Gas.Price

	for i := 0; i < Configuration.Transactions.Count; i++ {
		txResult, err := sendTransaction(currentNonce, gasPrice)

		if err != nil {
			fmt.Println(fmt.Sprintf("Error occurred: %s", err.Error()))
			return nil, err
		}

		txResults = append(txResults, txResult)

		currentNonce++
	}

	return txResults, nil
}

// AsyncBulkSendTransactions - sends transactions in bulk asynchronously
func AsyncBulkSendTransactions() {
	pools := 1

	currentNonce := Configuration.Account.Nonce
	gasPrice := Configuration.Transactions.Gas.Price
	fmt.Println(fmt.Sprintf("Running using network %s in %s mode", Configuration.Network.Name, strings.ToUpper(Configuration.Network.Mode)))
	fmt.Println(fmt.Sprintf("Current nonce is: %d, current gas price is %f", currentNonce, gasPrice))

	if Configuration.Transactions.Count > Configuration.Transactions.PoolSize {
		pools = int(math.RoundToEven(float64(Configuration.Transactions.Count) / float64(Configuration.Transactions.PoolSize)))
		fmt.Println(fmt.Sprintf("Number of goroutine pools: %d", pools))
	}

	for poolIndex := 0; poolIndex < pools; poolIndex++ {
		var waitGroup sync.WaitGroup

		if poolIndex > 1 {
			currentNonce := sdkNetwork.CurrentNonce(Configuration.Network.RPC, Configuration.Transactions.From)
			gasPrice = transactions.BumpGasPrice(gasPrice)
			fmt.Println(fmt.Sprintf("Nonce refreshed! Nonce is now: %d, gas price is now: %f", currentNonce, gasPrice))
		}

		for i := 0; i < Configuration.Transactions.PoolSize; i++ {
			waitGroup.Add(1)
			go asyncSendTransaction(currentNonce, gasPrice, &waitGroup)
			currentNonce++
		}

		waitGroup.Wait()
	}
}

func sendTransaction(currentNonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	txResult, err := executeTransaction(currentNonce, gasPrice)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error occurred: %s", err.Error()))
		return nil, err
	}

	txHash := txResult["transactionHash"].(string)

	fmt.Println(fmt.Sprintf("Receipt hash: %s", txHash))

	return txResult, nil
}

func asyncSendTransaction(currentNonce uint64, gasPrice numeric.Dec, waitGroup *sync.WaitGroup) {
	txResult, err := executeTransaction(currentNonce, gasPrice)

	if err == nil {
		txHash := txResult["transactionHash"].(string)
		fmt.Println(fmt.Sprintf("Receipt hash: %s", txHash))
	} else {
		fmt.Println(fmt.Sprintf("Error occurred: %s", err))
	}

	defer waitGroup.Done()
}

func executeTransaction(currentNonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	toAddress := utils.RandomStringSliceItem(r, Configuration.Transactions.Receivers)
	txResult, err := SendTransaction(Configuration.Account.Keystore, Configuration.Account.Account, Configuration.Network.RPC, Configuration.Network.API.ChainID, Configuration.Transactions.From, Configuration.Transactions.FromShard, toAddress, Configuration.Transactions.ToShard, Configuration.Transactions.Amount, Configuration.Transactions.Gas.Limit, gasPrice, currentNonce, Configuration.Transactions.Data, Configuration.Account.Passphrase, Configuration.Network.Node, Configuration.Transactions.Timeout)

	return txResult, err
}
