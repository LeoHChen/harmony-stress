package cmd

// Persistent is a collection of global/persistent flags
var Persistent PersistentFlags

// PersistentFlags represents the persistent flags
type PersistentFlags struct {
	Network         string
	NetworkMode     string
	Node            string
	Path            string
	ApplicationMode string
	From            string
	Passphrase      string
	Infinite        bool
	Count           int
	PoolSize        int
	Timeout         int
	Verbose         bool
	VerboseGoSdk    bool
	PprofPort       int

	TxMode      string
	FromShardID int
	ToShardID   int
}
