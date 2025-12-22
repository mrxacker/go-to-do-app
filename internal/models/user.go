package models

type UserID int64

type User struct {
	ID           UserID `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}
