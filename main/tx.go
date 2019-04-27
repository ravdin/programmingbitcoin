package main

type Transaction struct {
	Version    uint32              `json:"version"`
	Inputs     []TransactionInput  `json:"inputs"`
	Outputs    []TransactionOutput `json:"outputs"`
	Locktime   int                 `json:"locktime"`
	Testnet    bool                `json:"testnet"`
	Passphrase string              `json:"passphrase"`
}

type TransactionInput struct {
	PreviousValue string `json:"prevtx"`
	PreviousIndex int    `json:"previnput"`
}

type TransactionOutput struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}
