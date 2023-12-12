package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Info         UserInfo
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	PasswordHash string
}

type UserInfo struct {
	Email string
	Name  string
}
