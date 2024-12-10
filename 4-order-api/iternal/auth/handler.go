package auth

import (
	"net/http"
	"orderApi/configs"
	"orderApi/pkg/jwt"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth", handler.Auth())
	router.HandleFunc("POST /auth/verifying", handler.Verifying())
}

func (handler *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.BodyDecode[AuthRequest](&w, r)
		if err != nil {
			return
		}

		newUser, errNewUser := handler.AuthService.AuthCheckPhone(body.Phone)
		if errNewUser != nil {
			return
		}

		reqData := AuthResponse{
			SessionId: newUser.SessionId,
			Code:      3245,
		}

		response.Json(w, reqData, http.StatusOK)
	}
}

func (handler *AuthHandler) Verifying() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.BodyDecode[VerifyingRequests](&w, r)
		if err != nil {
			return
		}
		res, errVerifying := handler.AuthService.AuthCheckVerifying(body)
		if errVerifying != nil {
			response.Json(w, "Не правильный телефон или код", http.StatusUnauthorized)
			return
		}
		token, errToken := jwt.NewJWT(handler.Config.Auth.Secret).Create(res.Phone)
		if errToken != nil {
			response.Json(w, errToken.Error(), http.StatusInternalServerError)
			return
		}
		data := VerifyingResponses{
			Token: token,
		}

		response.Json(w, data, http.StatusOK)
	}
}
