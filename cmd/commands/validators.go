package commands

import (
	"path/filepath"

	cmdConfig "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/staking/validators"
	validatorsCreate "github.com/SebastianJ/harmony-stress/staking/validators/create"
	validatorsEdit "github.com/SebastianJ/harmony-stress/staking/validators/edit"
	"github.com/spf13/cobra"
)

func init() {
	cmdCreateValidators := &cobra.Command{
		Use:   "create-validators",
		Short: "Stress test validator creation",
		Long:  "Stress tests the creation of validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stressTestCreateValidators(cmd)
		},
	}

	RootCmd.AddCommand(cmdCreateValidators)

	cmdEditValidators := &cobra.Command{
		Use:   "edit-validators",
		Short: "Stress test validator editing",
		Long:  "Stress tests the editing of validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stressTestEditValidators(cmd)
		},
	}

	RootCmd.AddCommand(cmdEditValidators)
}

func stressTestCreateValidators(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := validators.Configure(basePath, cmdConfig.Persistent); err != nil {
		return err
	}

	validatorsCreate.StressTestValidatorCreation()

	return nil
}

func stressTestEditValidators(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := validators.Configure(basePath, cmdConfig.Persistent); err != nil {
		return err
	}

	validatorsEdit.StressTestValidatorEditing()

	return nil
}
