package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/gorilla/mux"
)

// MiddlewareValidateAuth validates the request and calls next if ok
func (kostHandler *KostHandler) MiddlewareValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// Get a session (existing/new)
		session, err := kostHandler.store.Get(r, "session-name")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		// check the token from the session
		// if token available, get the token from the session
		if session.Values["token"] == nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// determine the cookie value that holds the token
		tokenString := session.Values["token"].(string)

		if tokenString != "" {

			// Initialize a new instance of claims
			claims := &data.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error while parsing the token with claims")
				}

				return []byte(data.MySigningKey), nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					rw.WriteHeader(http.StatusUnauthorized)
					data.ToJSON(&GenericError{Message: err.Error()}, rw)

					return
				}

				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}

			if token.Valid {

				// create a new token for the current use, with a renewed expiration time
				expirationTime := time.Now().Add(time.Second * 86400 * 7)
				claims.StandardClaims.ExpiresAt = expirationTime.Unix()
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString([]byte(data.MySigningKey))

				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					data.ToJSON(&GenericError{Message: err.Error()}, rw)

					return
				}

				// renew the token in the session
				session.Options.MaxAge = 86400 * 7
				session.Values["token"] = tokenString
				session.Values["userLoggedin"] = claims.Username
				session.Save(r, rw)

				next.ServeHTTP(rw, r)
			} else {
				rw.WriteHeader(http.StatusUnauthorized)
				data.ToJSON(&GenericError{Message: "Token invalid"}, rw)

				return
			}
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&GenericError{Message: "Token invalid"}, rw)

			return
		}
	})
}

// MiddlewareParseKostGetRequest parses the kost payload from the query parameter
func (kostHandler *KostHandler) MiddlewareParseKostGetRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: "Unable to convert id"}, rw)

			return
		}

		// create the kost instance
		kost := &entities.Kost{
			ID: uint(id),
		}

		// add the kost to the context
		ctx := context.WithValue(r.Context(), KeyKost{}, kost)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

// MiddlewareParseKostPostRequest parses the kost payload in the request body from json
func (kostHandler *KostHandler) MiddlewareParseKostPostRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// validate content type to be application/json
		rw.Header().Add("Content-Type", "application/json")

		// create the kost instance
		kost := &entities.Kost{}

		// parse the request body to the given instance
		err := data.FromJSON(kost, r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		// add the kost to the context
		ctx := context.WithValue(r.Context(), KeyKost{}, kost)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

// MiddlewareParseKostAdsPostRequest parses the kost ads payload in the request body from json
func (kostHandler *KostHandler) MiddlewareParseKostAdsPostRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// validate content type to be application/json
		rw.Header().Add("Content-Type", "application/json")

		// create the kostAds instance
		kostAds := &entities.KostAds{}

		// parse the request body to the given instance
		err := data.FromJSON(kostAds, r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		// add the kost to the context
		ctx := context.WithValue(r.Context(), KeyKostAds{}, kostAds)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

// MiddlewareParseApprovalRequest parses the approval payload in the request body from json
func (kostHandler *KostHandler) MiddlewareParseApprovalRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// validate content type to be application/json
		rw.Header().Add("Content-Type", "application/json")

		// create the approval instance
		approval := &entities.ApprovalKost{}

		// parse the request body to the given instance
		err := data.FromJSON(approval, r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		// add the approval to the context
		ctx := context.WithValue(r.Context(), KeyApproval{}, approval)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

// MiddlewareParseUserRequest parses the user additional info payload in the request body from json
func (kostHandler *KostHandler) MiddlewareParseUserRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// get the user additional info from the header
		latitude := r.FormValue("latitude")
		longitude := r.FormValue("longitude")

		// create the user instance
		user := &entities.User{
			Latitude:  latitude,
			Longitude: longitude,
		}

		// add the user to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
