package verify

import (
	"net/http"
	"validation/pkg/response"
)

func ErrorsRegister(w http.ResponseWriter, err error) {
	if err != nil {
		response.JsonResponse(w, err.Error(), http.StatusBadRequest)
	}
}
