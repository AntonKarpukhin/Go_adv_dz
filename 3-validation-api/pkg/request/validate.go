package request

import "github.com/go-playground/validator/v10"

func IsValid[T any](data T) error {
	validate := validator.New()
	err := validate.Struct(data)
	return err
}
