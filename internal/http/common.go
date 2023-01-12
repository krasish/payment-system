package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/sirupsen/logrus"
)

type ClaimsKeyType string

const ClaimsCtxKey = ClaimsKeyType("context-claims")

func securedHandler(jwtKey []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerExtractor := request.AuthorizationHeaderExtractor
		bearerToken, err := headerExtractor.ExtractToken(r)
		if err != nil {
			respondWithMessage(w, fmt.Sprintf("could not extract token from auth header: %s", err.Error()), http.StatusBadRequest)
			return
		}
		//token claims validation is intentionally skipped for simplicity
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return jwtKey, nil
		})
		if err != nil {
			logrus.Errorf("could not parse token: %s", err.Error())
			respondWithMessage(w, "could not parse token", http.StatusUnauthorized)
			return
		}
		if token.Valid {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ClaimsCtxKey, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			respondWithMessage(w, "invalid token token", http.StatusUnauthorized)
		}
	}
}

func respondWithMessage(writer http.ResponseWriter, message string, statusCode int) {
	response := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to construct message response: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if _, err := writer.Write(resp); err != nil {
		logrus.Warnf("Failed to write response body: %v", err)
	}
}
