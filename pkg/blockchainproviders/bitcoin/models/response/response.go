package response

type GetBlockHashResponse struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

type GetBlockHeaderResponse struct {
	Result struct {
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
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    string      `json:"id"`
}

type GetBlockchainInfoResponse struct {
	Result struct {
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
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    string      `json:"id"`
}

type GetBlockCountResponse struct {
	Result int         `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}
