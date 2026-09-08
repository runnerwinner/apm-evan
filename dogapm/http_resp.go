package dogapm

import (
	"encoding/json"
	"net/http"
)
type Status struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Body any `json:"body"`
}

type httpStatus struct{}

var HttpStatus = &httpStatus{}

const (jsonContentType = "application/json")

func (h *httpStatus) OK(w http.ResponseWriter) {
	status := Status{
		Code:    http.StatusOK,
		Message: "OK",
	}
	data , _ := json.Marshal(status)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *httpStatus) OkBody(w http.ResponseWriter, msg string, body any) {
	status := Status{
		Code:    http.StatusOK,
		Message: msg,
		Body: body,
	}
	data , _ := json.Marshal(status)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *httpStatus) Fail(w http.ResponseWriter, msg string, body any) {
	status := Status{
		Code:    http.StatusBadRequest,
		Message: msg,
		Body:    body,
	}
	data , _ := json.Marshal(status)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(data)
}

func (h *httpStatus) Error(w http.ResponseWriter, msg string, body any) {
	status := Status{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Body:    body,
	}
	data , _ := json.Marshal(status)
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(data)
}