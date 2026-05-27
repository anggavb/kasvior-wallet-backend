package model

import "time"

type Receiver struct {
	Id          int
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
