package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/krasish/payment-system/internal/views"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/krasish/payment-system/internal/controllers"
)

type MerchantHandlerFactory struct {
	mc *controllers.MerchantController
	tc *controllers.TransactionController

	v *views.View
}

func NewMerchantHandlerFactory(mc *controllers.MerchantController, tc *controllers.TransactionController, v *views.View) *MerchantHandlerFactory {
	return &MerchantHandlerFactory{mc: mc, tc: tc, v: v}
}

func (f *MerchantHandlerFactory) BuildGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchants, err := f.mc.GetMerchants(r.Context())
		if err != nil {
			respondWithMessage(w, fmt.Sprintf("failed to get merchants: %v", err), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, merchants)
	}
}

func (f *MerchantHandlerFactory) BuildUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &controllers.Merchant{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
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
		if !strings.EqualFold(claims.Subject, m.Email) {
			respondWithMessage(w, "cannot update other merchants", http.StatusBadRequest)
			return
		}

		err := f.mc.UpdateMerchant(r.Context(), m)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to update merchant: %v", err)
			logrus.WithError(err).Error(errMsg)
			respondWithMessage(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
}

func (f *MerchantHandlerFactory) BuildDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
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
		if !strings.EqualFold(claims.Subject, email) {
			respondWithMessage(w, "cannot delete other merchants", http.StatusBadRequest)
			return
		}

		err := f.mc.DeleteMerchant(r.Context(), email)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to delete merchant: %v", err)
			logrus.WithError(err).Error(errMsg)
			respondWithMessage(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
}

func (f *MerchantHandlerFactory) BuildHTMLTemplateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchants, err := f.mc.GetMerchants(r.Context())
		if err != nil {
			respondWithMessage(w, "cannot get merchants", http.StatusInternalServerError)
		}
		transactions, err := f.tc.GetTransactions(r.Context())
		if err != nil {
			respondWithMessage(w, "cannot get merchants", http.StatusInternalServerError)
		}
		viewData := views.NewMerchantsData(merchants, transactions)

		if err := f.v.Render(w, viewData); err != nil {
			respondWithMessage(w, "cannot parse HTML template", http.StatusInternalServerError)
		}
	}
}
