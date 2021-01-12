package handlers

import (
	"net/http"
)

// AddKost is a method to add the new given kost info
func (kostHandler *KostHandler) AddKost(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}
