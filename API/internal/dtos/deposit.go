package dtos

type CreateUserStaticPixDepositInput struct {
	WalletAddress string `json:"walletAddress"`
	Chain         string `json:"chain"`
	Amount        int    `json:"amount"`
}

type CreateStaticPixDepositOutput struct {
	Id string `json:"id"`
}
