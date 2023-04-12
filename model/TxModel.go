package model

import "projectx/crypto"

type TransactionCreate struct {
	Provider string   `json:"provider"`
	Track    []string `json:"track"`
	PrivKey  string   `json:"privkey"`
}

type TransactionResponse struct {
	Provider  string            `json:"provider"`
	Track     []string          `json:"track"`
	Signature *crypto.Signature `json:"Signature"`
}
