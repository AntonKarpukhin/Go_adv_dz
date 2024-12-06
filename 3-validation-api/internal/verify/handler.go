package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"validation/pkg/account"
	"validation/pkg/file"
	"validation/pkg/request"
	"validation/pkg/utils"
)

type VerifierHandler struct {
	Accounts []*account.Account
	Db       *file.JsFile
	mu       sync.Mutex
}

func NewVerifierHandler(router *http.ServeMux, db *file.JsFile) {
	handler := &VerifierHandler{
		Accounts: make([]*account.Account, 0),
		Db:       db,
	}

	router.HandleFunc("POST /auth/verify", handler.Register())
	router.HandleFunc("/auth/verify/{hash}", handler.CheckEmail())
}

func (verifier *VerifierHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		ErrorsRegister(w, err)
		err = request.IsValid(payload)
		ErrorsRegister(w, err)

		newAccount := account.NewAccount(payload.Email)
		verifier.mu.Lock()
		verifier.Accounts = append(verifier.Accounts, newAccount)
		verifier.mu.Unlock()

		dbData, _ := verifier.Db.Read()
		var accountsFromFile []*account.Account
		err = json.Unmarshal(dbData, &accountsFromFile)
		if err != nil {
			ErrorsRegister(w, errors.New("Ошибка при обработке JSON:"))
			return
		}

		for _, dataAccount := range accountsFromFile {
			isMatch := utils.IsEmail(dataAccount.Email, payload.Email)
			if isMatch {
				ErrorsRegister(w, errors.New("такой аккаунт уже существует"))
				return
			}
		}

		marshal, _ := json.Marshal(verifier.Accounts)
		verifier.Db.Write(marshal)

		CreateMail(newAccount.Password, w)
	}
}

func (verifier *VerifierHandler) CheckEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		dbData, _ := verifier.Db.Read()
		var accountsFromFile []*account.Account
		err := json.Unmarshal(dbData, &accountsFromFile)
		if err != nil {
			ErrorsRegister(w, errors.New("Ошибка при обработке JSON:"))
			return
		}

		var newAccounts []account.Account

		for _, dataAccount := range accountsFromFile {
			isMatch := utils.IsEmail(dataAccount.Password, hash)
			if isMatch {
				newAccounts = append(newAccounts, *dataAccount)
			} else {
				fmt.Println("Аккаунт удален")
			}
		}

		marshal, _ := json.Marshal(newAccounts)
		verifier.Db.Write(marshal)
	}
}
