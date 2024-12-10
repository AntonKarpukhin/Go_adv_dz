package auth

import (
	"math/rand"
	"orderApi/iternal/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-*!")

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) AuthCheckPhone(phone string) (*user.User, error) {
	userData, err := service.UserRepository.FindByPhone(phone)
	if err != nil {
		var newUser user.User
		newUser.Phone = phone
		newUser.Code = 3245
		newUser.SessionId = service.GenerateSessionId()
		_, err = service.UserRepository.Create(&newUser)
		if err != nil {
			return nil, err
		}
		return &newUser, nil
	}
	userData.SessionId = service.GenerateSessionId()

	err = service.UserRepository.UpdateSessionId(userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (service *AuthService) AuthCheckVerifying(user *VerifyingRequests) (*user.User, error) {
	answer, err := service.UserRepository.FindBySessionId(user.SessionId, user.Code)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (service *AuthService) GenerateSessionId() string {
	res := make([]rune, 10)
	for i := range res {
		res[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(res)
}
