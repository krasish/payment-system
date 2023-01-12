package http

import (
	"net/http"
	"os"

	"github.com/krasish/payment-system/internal/views"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/krasish/payment-system/internal/config"
	"github.com/krasish/payment-system/internal/controllers"
)

func CreateHTTPServer(cfg config.HttpConfig, tc *controllers.TransactionController, mc *controllers.MerchantController, v *views.View) (*http.Server, error) {
	var (
		mainRouter = mux.NewRouter()
		jwtKey     = []byte(cfg.JwtKey)
	)

	//Transaction handlers
	transactionHandlerFactory := NewTransactionHandlerFactory(tc)

	getTransactionHandler := transactionHandlerFactory.BuildGetHandler()
	createTransactionHandler := securedHandler(jwtKey, handlers.ContentTypeHandler(transactionHandlerFactory.BuildCreateHandler(), ContentTypeAppJSON).ServeHTTP)

	mainRouter.HandleFunc(cfg.TransactionPath, getTransactionHandler).Methods(http.MethodGet)
	mainRouter.HandleFunc(cfg.TransactionPath, createTransactionHandler).Methods(http.MethodPost)

	//Merchant handlers
	merchantHandlerFactory := NewMerchantHandlerFactory(mc, tc, v)

	getMerchantHandler := merchantHandlerFactory.BuildGetHandler()
	updateMerchantHandler := securedHandler(jwtKey, handlers.ContentTypeHandler(merchantHandlerFactory.BuildUpdateHandler(), ContentTypeAppJSON).ServeHTTP)
	deleteMerchantHandler := securedHandler(jwtKey, merchantHandlerFactory.BuildDeleteHandler())
	htmlTemplateHandler := merchantHandlerFactory.BuildHTMLTemplateHandler()

	mainRouter.HandleFunc(cfg.MerchantPath, getMerchantHandler).Methods(http.MethodGet)
	mainRouter.HandleFunc(cfg.MerchantPath, updateMerchantHandler).Methods(http.MethodPut)
	mainRouter.HandleFunc(cfg.MerchantPath, deleteMerchantHandler).Methods(http.MethodDelete)

	viewsRouter := mainRouter.PathPrefix(cfg.ViewsPath).Subrouter()
	viewsRouter.HandleFunc(cfg.MerchantPath, htmlTemplateHandler)

	loggingHandler := handlers.LoggingHandler(os.Stdout, handlers.RecoveryHandler()(mainRouter))
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           loggingHandler,
		ReadHeaderTimeout: cfg.ServerTimeout,
	}

	return srv, nil
}
