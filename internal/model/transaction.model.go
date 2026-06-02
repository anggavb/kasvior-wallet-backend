package model

import "time"

type Receiver struct {
	Id          int
	WalletId    string
	Photo       *string
	Receiver    *string
	PhoneNumber *string
}

type Transaction struct {
	Id        int
	WalletId  string
	Amount    int
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type TransferTransaction struct {
	Id                int
	SenderWalletId    string
	RecipientWalletId string
	Amount            float64
	Status            string
}

type TransactionHistoryItem struct {
	Id                int
	Type              string
	Direction         string
	Status            string
	Amount            float64
	CounterpartyName  *string
	CounterpartyPhone *string
	CounterpartyPhoto *string
	PaymentMethod     *string
	Notes             *string
	CreatedAt         time.Time
}
