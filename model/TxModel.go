package model

import "projectx/crypto"

type TransactionCreate struct {
	Provider string   `json:"provider"`
	Track    []string `json:"track"`
}

type TransactionResponse struct {
	Provider  string   `json:"provider"`
	Track     []string `json:"track"`
	From      crypto.PublicKey
	Signature *crypto.Signature
}
