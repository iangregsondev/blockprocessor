package request

type GetBlockOptions struct {
	Encoding                       string
	TransactionDetails             string
	Rewards                        bool
	MaxSupportedTransactionVersion int
}

type GetTransactionOptions struct {
	Commitment                     string
	Encoding                       string
	MaxSupportedTransactionVersion int
}
