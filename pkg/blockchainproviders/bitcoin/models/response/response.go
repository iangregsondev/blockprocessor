package response

type GetBlockHashResponse struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

type GetBlockHeaderResponse struct {
	Result BlockHeader `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

type BlockHeader struct {
	Hash              string  `json:"hash"`
	Confirmations     int     `json:"confirmations"`
	Height            int     `json:"height"`
	Version           int     `json:"version"`
	VersionHex        string  `json:"versionHex"`
	Merkleroot        string  `json:"merkleroot"`
	Time              int     `json:"time"`
	Mediantime        int     `json:"mediantime"`
	Nonce             int64   `json:"nonce"`
	Bits              string  `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
	Chainwork         string  `json:"chainwork"`
	NTx               int     `json:"nTx"`
	Previousblockhash string  `json:"previousblockhash"`
	Nextblockhash     string  `json:"nextblockhash"`
}

type GetBlockchainInfoResponse struct {
	Result BlockchainInfo `json:"result"`
	Error  interface{}    `json:"error"`
	ID     string         `json:"id"`
}

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Time                 int     `json:"time"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Initialblockdownload bool    `json:"initialblockdownload"`
	Chainwork            string  `json:"chainwork"`
	SizeOnDisk           int64   `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Warnings             string  `json:"warnings"`
}

type GetBlockCountResponse struct {
	Result int         `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

type GetBlockResponse struct {
	Result Block       `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

type Block struct {
	Hash              string   `json:"hash"`
	Confirmations     int      `json:"confirmations"`
	Height            int      `json:"height"`
	Version           int      `json:"version"`
	VersionHex        string   `json:"versionHex"`
	Merkleroot        string   `json:"merkleroot"`
	Time              int      `json:"time"`
	Mediantime        int      `json:"mediantime"`
	Nonce             int      `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	NTx               int      `json:"nTx"`
	Previousblockhash string   `json:"previousblockhash"`
	Nextblockhash     string   `json:"nextblockhash"`
	Strippedsize      int      `json:"strippedsize"`
	Size              int      `json:"size"`
	Weight            int      `json:"weight"`
	Tx                []string `json:"tx"`
}

type GetRawTransactionResponse struct {
	Result RawTransaction `json:"result"`
	Error  interface{}    `json:"error"`
	ID     string         `json:"id"`
}

type RawTransaction struct {
	Txid          string `json:"txid"`
	Hash          string `json:"hash"`
	Version       int    `json:"version"`
	Size          int    `json:"size"`
	Vsize         int    `json:"vsize"`
	Weight        int    `json:"weight"`
	Locktime      int    `json:"locktime"`
	Vin           []VIn  `json:"vin"`
	Vout          []VOut `json:"vout"`
	Hex           string `json:"hex"`
	Blockhash     string `json:"blockhash"`
	Confirmations int    `json:"confirmations"`
	Time          int    `json:"time"`
	Blocktime     int    `json:"blocktime"`
}

type VIn struct {
	Coinbase    string   `json:"coinbase"`
	Txinwitness []string `json:"txinwitness"`
	Sequence    int64    `json:"sequence"`
}

type VOut struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
