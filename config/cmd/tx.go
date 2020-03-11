package cmd

// Tx is a collection of tx related options
var Tx TxFlags

// TxFlags represents the tx flags
type TxFlags struct {
	Mode        string
	FromShardID int
	ToShardID   int
	Amount      string
}
