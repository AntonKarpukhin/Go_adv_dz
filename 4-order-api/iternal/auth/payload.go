package auth

type AuthRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

type AuthResponse struct {
	SessionId string `json:"sessionId"`
	Code      int    `json:"code"`
}

type VerifyingRequests struct {
	SessionId string `json:"sessionId" required:"true"`
	Code      int    `json:"code" required:"true"`
}

type VerifyingResponses struct {
	Token string `json:"token"`
}
