package request

import (
	"finance-tracker/pkg/response"
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)

	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	err = isValid(body)
	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
