package models

type TransactionPartnerUser struct {
	Type              string        `json:"type"`             
	TransactionType   string        `json:"transaction_type"`  
	User              TelegramUser  `json:"user"`              
}

type TelegramUser struct {
	ID       int64  `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Username string `json:"username"`
}

type StarTransaction struct {
	ID             string                `json:"id"`
	Amount         int                   `json:"amount"`
	NanostarAmount int                   `json:"nanostar_amount"`
	Date           int64                 `json:"date"`
	Source         TransactionPartnerUser `json:"source"`
	Receiver       TransactionPartnerUser `json:"receiver"`
}