package entity

type Wallet struct {
	id               int64
	userId           int64
	TotalBalance     int64 `json:"total_balance,omitempty"`
	AvailableBalance int64 `json:"available_balance,omitempty"`
}

func NewWallet(id int64, userId int64, totalBalance int64, availableBalance int64) *Wallet {
	return &Wallet{
		id:               id,
		userId:           userId,
		TotalBalance:     totalBalance,
		AvailableBalance: availableBalance,
	}
}
