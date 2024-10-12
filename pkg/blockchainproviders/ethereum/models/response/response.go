package response

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BlockNumberResponse struct {
	Result  string `json:"result"`
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Error   *Error `json:"error,omitempty"`
}

type BlockByNumberResponse struct {
	Result  BlockByNumber `json:"result"`
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Error   *Error        `json:"error,omitempty"`
}

type BlockByNumber struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []string      `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
}
