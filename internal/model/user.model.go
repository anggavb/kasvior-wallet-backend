package model

import "time"

type User struct {
	Id          int        `db:"id"`
	Fullname    *string    `db:"fullname"`
	Email       string     `db:"email"`
	Password    string     `db:"password"`
	Pin         *string    `db:"pin"`
	PhoneNumber *string    `db:"phone_number"`
	Photo       *string    `db:"photo"`
	IsVerified  bool       `db:"is_verified"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type UserDashboardInformation struct {
	Balance float64 `db:"balance"`
	Income  float64 `db:"income"`
	Expense float64 `db:"expense"`
}
