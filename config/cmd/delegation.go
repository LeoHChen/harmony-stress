package cmd

// Delegation is a collection of delegation related options
var Delegation DelegationFlags

// DelegationFlags represents the delegation flags
type DelegationFlags struct {
	Mode        string
	FromShardID int
	ToShardID   int
	Amount      string
	GasPrice    string
}
