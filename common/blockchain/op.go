package blockchain

const (
	// PricingOP the account using for pricing
	PricingOP = "pricingOP"
	// DepositOP the account using for deposit to exchange
	DepositOP = "depositOP"
)

// MinedNonePicker just an interface container shared function of core/blockchain and fetcher/blockchain interface
type MinedNonePicker interface {
	GetMinedNonceWithOP(op string) (uint64, error)
}
