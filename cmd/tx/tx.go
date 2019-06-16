package main

type transaction struct {
	Version    uint32              `json:"version"`
	Inputs     []transactionInput  `json:"inputs"`
	Outputs    []transactionOutput `json:"outputs"`
	Locktime   int                 `json:"locktime"`
	Testnet    bool                `json:"testnet"`
	Passphrase string              `json:"passphrase"`
}

type transactionInput struct {
	PreviousValue string `json:"prevtx"`
	PreviousIndex int    `json:"previnput"`
}

type transactionOutput struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}
