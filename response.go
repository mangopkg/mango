package mango

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Error      bool        `json:"error"`
}

func (r *Response) Send(w http.ResponseWriter) {
	w.WriteHeader(r.StatusCode)
	if r.Error {
		r.Data = nil
		if err := json.NewEncoder(w).Encode(r); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(r); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
