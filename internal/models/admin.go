package models

import "5g-v2x-user-service/internal/utils"

type Admin struct {
	Username       string `bson:"username"`
	HashedPassword string `bson:"hashed_password"`
}

// VerifyPassword is ...
func (a *Admin) VerifyPassword(pwd string) bool {
	hashed, err := utils.HashAndSalt([]byte(pwd))
	if err != nil {
		return false
	}
	if hashed != a.HashedPassword {
		return false
	}
	return true
}
