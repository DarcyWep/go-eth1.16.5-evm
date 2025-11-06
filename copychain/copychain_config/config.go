package copychain_config

import "math/big"

var (
	StartBlockNumber  *big.Int = new(big.Int).SetInt64(21_000_000)
	FinishBlockNumber *big.Int = new(big.Int).SetInt64(23_700_001)
	AddSpan           *big.Int = new(big.Int).SetInt64(1)
)
