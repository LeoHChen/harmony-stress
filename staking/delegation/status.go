package delegation

import (
	"fmt"
	"math"
	"sync"

	sdkValidator "github.com/SebastianJ/harmony-sdk/staking/validator"
)

// DelegationStatistics - lookup active delegations on the network
func DelegationStatistics(onlyActive bool) error {
	fmt.Printf("Looking up delegation statistics - network: %s, mode: %s, node: %s\n", Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Node)

	pageSize := 10
	pages := 1
	totalCount := 0
	var validators []string
	var err error
	message := ""

	if onlyActive {
		message = " elected"
		validators, err = sdkValidator.AllElected(Configuration.Network.RPC)
	} else {
		validators, err = sdkValidator.All(Configuration.Network.RPC)
	}

	if err != nil {
		return err
	}

	validatorCount := len(validators)
	fmt.Println(fmt.Sprintf("Found a total of %d%s validators to look up delegation statistics for", validatorCount, message))
	pages = calculatePageCount(validatorCount, pageSize)
	totalChecked := 0

	delegationsCountChannel := make(chan int, validatorCount)
	var waitGroup sync.WaitGroup

	for page := 0; page < pages; page++ {
		for i := 0; i < pageSize; i++ {
			position, ok := processable(page, pageSize, i, validatorCount)
			if ok {
				waitGroup.Add(1)
				address := validators[position]
				go lookupValidator(address, delegationsCountChannel, &waitGroup)
				totalChecked++
			}
		}

		waitGroup.Wait()
	}

	close(delegationsCountChannel)

	for delegationCount := range delegationsCountChannel {
		totalCount = totalCount + delegationCount
	}

	fmt.Printf("Total number of delegations on the network is: %d - network: %s, mode: %s, node: %s\n", totalCount, Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Node)
	fmt.Printf("Total checked number of validators: %d\n", totalChecked)

	return nil
}

func calculatePageCount(totalCount int, pageSize int) int {
	if totalCount > 0 {
		pageNumber := math.RoundToEven(float64(totalCount) / float64(pageSize))
		if math.Mod(float64(totalCount), float64(pageSize)) > 0 {
			return int(pageNumber) + 1
		}

		return int(pageNumber)
	} else {
		return 0
	}
}

func processable(page int, pageSize int, index int, totalCount int) (position int, ok bool) {
	position = ((page * pageSize) + index)
	ok = position <= (totalCount - 1)
	return position, ok
}

func lookupValidator(address string, delegationsCountChannel chan<- int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	result, err := sdkValidator.Information(Configuration.Network.Node, address)

	if err == nil {
		validator := result.Validator
		delegationsCount := len(validator.Delegations)
		fmt.Printf("Delegation count for validator %s is: %d\n", address, delegationsCount)
		delegationsCountChannel <- delegationsCount
	}
}
