package api

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	object interface{}
}

func SendJson(object interface{}) JsonResponse {
	return JsonResponse{
		object: object,
	}
}

func (result JsonResponse) Write(w *http.ResponseWriter) error {
	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
	(*w).WriteHeader(http.StatusOK)

	jsonBytes, err := json.Marshal(result.object)
	if err != nil {
		return err
	}

	_, err = (*w).Write(jsonBytes)
	return err
}
