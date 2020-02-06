package configuration

// EthereumNodeConfiguration ...
type EthereumNodeConfiguration struct {
	Main   string
	Backup []string
}

// NewEthereumNodeConfiguration returns a new Ethereum node configuration.
func NewEthereumNodeConfiguration(main string, backup []string) *EthereumNodeConfiguration {
	return &EthereumNodeConfiguration{Main: main, Backup: backup}
}
