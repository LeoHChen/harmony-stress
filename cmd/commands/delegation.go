package commands

import (
	"path/filepath"
	"strings"

	cmdConfig "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/staking/delegation"
	delegationGeneration "github.com/SebastianJ/harmony-stress/staking/delegation/generate"
	"github.com/spf13/cobra"
)

func init() {
	cmdDelegation := &cobra.Command{
		Use:   "delegations",
		Short: "Stress test delegation transactions",
		Long:  "Stress test delegation transactions (delegate / undelegate)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stressTestDelegations(cmd)
		},
	}

	cmdConfig.Delegation = cmdConfig.DelegationFlags{}
	cmdDelegation.Flags().StringVar(&cmdConfig.Delegation.ValidatorAddress, "validator-address", "", "<amount>")
	cmdDelegation.Flags().IntVar(&cmdConfig.Delegation.FromShardID, "from-shard", 0, "<shardID>")
	cmdDelegation.Flags().IntVar(&cmdConfig.Delegation.ToShardID, "to-shard", 0, "<shardID>")
	cmdDelegation.Flags().StringVar(&cmdConfig.Delegation.Amount, "amount", "", "<amount>")
	cmdDelegation.Flags().BoolVar(&cmdConfig.Delegation.OnlyActive, "active", false, "<active>")
	cmdDelegation.Flags().StringVar(&cmdConfig.Delegation.DelegationMode, "delegation-mode", "generate", "<delegation-mode>")
	RootCmd.AddCommand(cmdDelegation)

	cmdDelegationStatistics := &cobra.Command{
		Use:   "delegations-info",
		Short: "Check info for delegations",
		Long:  "Check info for delegations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return delegationStatistics(cmd)
		},
	}
	cmdDelegationStatistics.Flags().BoolVar(&cmdConfig.Delegation.OnlyActive, "active", false, "<active>")
	RootCmd.AddCommand(cmdDelegationStatistics)
}

func stressTestDelegations(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := delegation.Configure(basePath, cmdConfig.Persistent, cmdConfig.Delegation); err != nil {
		return err
	}

	if strings.ToLower(cmdConfig.Delegation.DelegationMode) == "generate" {
		if err := delegationGeneration.StressTestDelegations(); err != nil {
			return err
		}
	} else {
		if _, err := delegation.StressTestDelegations(cmdConfig.Delegation.ValidatorAddress, cmdConfig.Delegation.OnlyActive); err != nil {
			return err
		}
	}

	return nil
}
func delegationStatistics(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := delegation.Configure(basePath, cmdConfig.Persistent, cmdConfig.Delegation); err != nil {
		return err
	}

	if err := delegation.DelegationStatistics(cmdConfig.Delegation.OnlyActive); err != nil {
		return err
	}

	return nil
}
