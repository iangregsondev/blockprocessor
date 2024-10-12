package response

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BlockHeightResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  int64  `json:"result"`
	ID      string `json:"id"`
	Error   *Error `json:"error,omitempty"`
}

type BlockResponse struct {
	Result  Block  `json:"result"`
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Error   *Error `json:"error,omitempty"`
}

type Block struct {
	BlockHeight       interface{} `json:"blockHeight"`
	BlockTime         interface{} `json:"blockTime"`
	Blockhash         string      `json:"blockhash"`
	ParentSlot        int         `json:"parentSlot"`
	PreviousBlockhash string      `json:"previousBlockhash"`
	Signatures        []string    `json:"signatures"`
}

type TransactionResponse struct {
	Result  Transaction `json:"result"`
	Jsonrpc string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Error   *Error      `json:"error,omitempty"`
}

type Transaction struct {
	BlockTime   interface{}      `json:"blockTime"`
	Meta        Meta             `json:"meta"`
	Slot        int              `json:"slot"`
	Transaction InnerTransaction `json:"transaction"`
	Version     string           `json:"version"`
}

type InnerTransaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

type Message struct {
	AccountKeys     []string             `json:"accountKeys"`
	Header          MessageHeader        `json:"header"`
	Instructions    []MessageInstruction `json:"instructions"`
	RecentBlockhash string               `json:"recentBlockhash"`
}

type MessageHeader struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}

type MessageInstruction struct {
	Accounts       []int       `json:"accounts"`
	Data           string      `json:"data"`
	ProgramIDIndex int         `json:"programIdIndex"`
	StackHeight    interface{} `json:"stackHeight"`
}

type Meta struct {
	Err               interface{}   `json:"err"`
	Fee               int           `json:"fee"`
	InnerInstructions interface{}   `json:"innerInstructions"`
	LoadedAddresses   LoadedAddress `json:"loadedAddresses"`
	LogMessages       interface{}   `json:"logMessages"`
	PostBalances      []interface{} `json:"postBalances"`
	PostTokenBalances interface{}   `json:"postTokenBalances"`
	PreBalances       []interface{} `json:"preBalances"`
	PreTokenBalances  interface{}   `json:"preTokenBalances"`
	Rewards           interface{}   `json:"rewards"`
	Status            MetaStatus    `json:"status"`
}

type LoadedAddress struct {
	Readonly []interface{} `json:"readonly"`
	Writable []interface{} `json:"writable"`
}

type MetaStatus struct {
	Ok interface{} `json:"Ok"`
}
