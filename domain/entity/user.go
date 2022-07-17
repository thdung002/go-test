package entity

type User struct {
	Address    []byte `json:"address"`
	Amount     []byte `json:"amount"`
	Nonce      uint64 `json:"nonce"`
	PrivateKey []byte `json:"private_key"`
}
