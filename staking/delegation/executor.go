package delegation

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	sdkValidator "github.com/SebastianJ/harmony-sdk/staking/validator"
	"github.com/SebastianJ/harmony-sdk/transactions"
	"github.com/SebastianJ/harmony-stress/utils"
	"github.com/harmony-one/harmony/numeric"
)

// StressTestDelegations - stress tests delegation functionality
func StressTestDelegations() (txResults []map[string]interface{}, err error) {
	pools := 1

	currentNonce := Configuration.Account.Nonce
	gasPrice := Configuration.Delegation.Gas.Price
	fmt.Println(fmt.Sprintf("Running using network %s in %s mode", Configuration.Network.Name, strings.ToUpper(Configuration.Network.Mode)))
	fmt.Println(fmt.Sprintf("Current nonce is: %d, current gas price is %f", currentNonce, gasPrice))

	electedValidators, err := sdkValidator.AllElected(Configuration.Network.RPC)
	fmt.Println(fmt.Sprintf("Found a total of %d elected validators to send delegations to", len(electedValidators)))

	if Configuration.Application.Count > Configuration.Application.PoolSize {
		pools = int(math.RoundToEven(float64(Configuration.Application.Count) / float64(Configuration.Application.PoolSize)))
		fmt.Println(fmt.Sprintf("Number of goroutine pools: %d", pools))
	}

	for poolIndex := 0; poolIndex < pools; poolIndex++ {
		var waitGroup sync.WaitGroup

		if poolIndex > 1 {
			currentNonce = sdkNetwork.CurrentNonce(Configuration.Network.RPC, Configuration.Application.From)
			gasPrice = transactions.BumpGasPrice(gasPrice)
			fmt.Println(fmt.Sprintf("Nonce refreshed! Nonce is now: %d, gas price is now: %f", currentNonce, gasPrice))
		}

		for i := 0; i < Configuration.Application.PoolSize; i++ {
			r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
			toAddress := utils.RandomStringSliceItem(r, electedValidators)

			if Configuration.Application.Mode == "sync" {
				txResult, err := sendDelegation(toAddress, currentNonce, gasPrice)
				if err != nil {
					fmt.Println(fmt.Sprintf("Error occurred: %s", err.Error()))
					return nil, err
				}

				txResults = append(txResults, txResult)
			} else if Configuration.Application.Mode == "async" {
				waitGroup.Add(1)
				go asyncSendDelegation(toAddress, currentNonce, gasPrice, &waitGroup)
			}

			currentNonce = currentNonce + 1
		}

		if Configuration.Application.Mode == "async" {
			waitGroup.Wait()
		}
	}

	return txResults, nil
}

func sendDelegation(address string, currentNonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	return executeDelegation(address, currentNonce, gasPrice)
}

func asyncSendDelegation(address string, currentNonce uint64, gasPrice numeric.Dec, waitGroup *sync.WaitGroup) {
	executeDelegation(address, currentNonce, gasPrice)
	defer waitGroup.Done()
}

func executeDelegation(address string, currentNonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	txResult, err := Delegate(address, currentNonce, gasPrice)
	if err == nil {
		txHash := txResult["transactionHash"].(string)
		fmt.Println(fmt.Sprintf("Sent delegation of %f from %s to %s, nonce: %d, gas price: %f, tx hash: %s", Configuration.Delegation.Amount, Configuration.Account.Account.Address, address, currentNonce, gasPrice, txHash))
	} else {
		fmt.Println(fmt.Sprintf("Error occurred while sending delegation of %f from %s to %s, nonce: %d - error: %s", Configuration.Delegation.Amount, Configuration.Account.Account.Address, address, currentNonce, err.Error()))
	}

	return txResult, err
}
