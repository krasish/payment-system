package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/sirupsen/logrus"

	"github.com/krasish/payment-system/internal/controllers"
)

type TransactionHandlerFactory struct {
	tc *controllers.TransactionController
}

func NewTransactionHandlerFactory(tc *controllers.TransactionController) *TransactionHandlerFactory {
	return &TransactionHandlerFactory{tc: tc}
}

func (f *TransactionHandlerFactory) BuildCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := &controllers.Transaction{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			logrus.WithError(err).Error("Failed to read transaction from request body")
			respondWithMessage(w, "could not read request body", http.StatusInternalServerError)
			return
		}
		value := r.Context().Value(ClaimsCtxKey)
		if value == nil {
			respondWithMessage(w, "could not get token claims", http.StatusUnauthorized)
			return
		}
		claims, ok := value.(jwt.StandardClaims)
		if !ok {
			respondWithMessage(w, "invalid token claims format", http.StatusUnauthorized)
			return
		}
		if !strings.EqualFold(claims.Subject, t.MerchantEmail) {
			respondWithMessage(w, "cannot create transactions on behalf of other merchants", http.StatusBadRequest)
			return
		}

		err := f.tc.CreateTransaction(r.Context(), t)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to create transaction: %v", err)
			logrus.WithError(err).Error(errMsg)
			respondWithMessage(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
}

func (f *TransactionHandlerFactory) BuildGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := f.tc.GetTransactions(r.Context())
		if err != nil {
			respondWithMessage(w, fmt.Sprintf("failed to get transactions: %v", err), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, transactions)
	}
}
