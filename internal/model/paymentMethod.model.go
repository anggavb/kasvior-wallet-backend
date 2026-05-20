package model

import "time"

type PaymentMethod struct {
	Id        int        `db:"id"`
	Name      string     `db:"name"`
	Logo      string     `db:"logo"`
	Method    string     `db:"method"`
	Tax       int        `db:"tax"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
