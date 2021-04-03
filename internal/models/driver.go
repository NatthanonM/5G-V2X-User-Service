package models

import (
	"5g-v2x-user-service/internal/utils"
	"time"
)

type Driver struct {
	DriverID       string     `bson:"_id"`
	Firstname      *string    `bson:"firstname"`
	Lastname       *string    `bson:"lastname"`
	Username       string     `bson:"username"`
	Password       string     `bson:"-"`
	HashedPassword string     `bson:"hashed_password"`
	DateOfBirth    *time.Time `bson:"date_of_birth"`
	Gender         string     `bson:"gender"`
	DeletedAt      *time.Time `bson:"deleted_at"`
}

// VerifyPassword is ...
func (d *Driver) VerifyPassword(pwd string) bool {
	hashed, err := utils.HashAndSalt([]byte(pwd))
	if err != nil {
		return false
	}
	if hashed != d.HashedPassword {
		return false
	}
	return true
}
