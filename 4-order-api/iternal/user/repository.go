package user

import (
	"errors"
	"orderApi/pkg/db"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByPhone(phone string) (*User, error) {
	var user User
	result := repo.Database.First(&user, "phone = ?", phone)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindBySessionId(sessionId string, code int) (*User, error) {
	var user User
	result := repo.Database.First(&user, "session_id = ?", sessionId)
	if result.Error != nil {
		return nil, result.Error
	}

	if user.Code != code {
		return nil, errors.New("code not match")
	}

	return &user, nil
}

func (repo *UserRepository) UpdateSessionId(user *User) error {
	result := repo.Database.DB.Model(&User{}).Where("phone = ?", user.Phone).Update("session_id", user.SessionId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
