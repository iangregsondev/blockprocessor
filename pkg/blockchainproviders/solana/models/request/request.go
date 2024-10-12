package request

type GetBlockOptions struct {
	Encoding                       string
	TransactionDetails             string
	Rewards                        bool
	MaxSupportedTransactionVersion int
}
