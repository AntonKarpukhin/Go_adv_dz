package account

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-*!")

type Account struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

func NewAccount(email string) *Account {
	newAccount := &Account{
		Email: email,
	}
	newAccount.GeneratePassword()
	return newAccount
}

func (a *Account) GeneratePassword() {
	res := make([]rune, 10)
	for i := range res {
		res[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	a.Password = string(res)
}
