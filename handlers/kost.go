package handlers

import (
	"github.com/fakhripraya/kost-service/data"

	"github.com/hashicorp/go-hclog"
	"github.com/srinathgs/mysqlstore"
)

// KeyKost is a key used for the Kost object in the context
type KeyKost struct{}

// KeyKostAds is a key used for the Kost Ads object in the context
type KeyKostAds struct{}

// KeyApproval is a key used for the Approval object in the context
type KeyApproval struct{}

// KeyUser is a key used for the User object in the context
type KeyUser struct{}

// KostHandler is a handler struct for kost changes
type KostHandler struct {
	logger hclog.Logger
	kost   *data.Kost
	store  *mysqlstore.MySQLStore
}

// NewKostHandler returns a new Kost handler with the given logger
func NewKostHandler(newLogger hclog.Logger, newKost *data.Kost, newStore *mysqlstore.MySQLStore) *KostHandler {
	return &KostHandler{newLogger, newKost, newStore}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}
