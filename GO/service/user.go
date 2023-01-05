package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName       string
	HashedPassword string
	Role           string
}

func NewUser(userName string, password string, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	return &User{
		UserName:       userName,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}, nil
}

func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

func (user *User) Clone() *User {
	return &User{
		UserName:       user.UserName,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}
