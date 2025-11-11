package entity

type Wallet struct {
	id               int64
	userId           int64
	totalBalance     int64
	availableBalance int64
}

func NewWallet(id int64, userId int64, totalBalance int64, availableBalance int64) *Wallet {
	return &Wallet{
		id:               id,
		userId:           userId,
		totalBalance:     totalBalance,
		availableBalance: availableBalance,
	}
}

func (w *Wallet) TotalBalance() int64 {
	return w.totalBalance
}

func (w *Wallet) AvailableBalance() int64 {
	return w.availableBalance
}
