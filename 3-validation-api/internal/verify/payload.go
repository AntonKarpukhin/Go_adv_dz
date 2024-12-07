package verify

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}
