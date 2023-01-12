package http

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/krasish/payment-system/internal/config"
	"github.com/krasish/payment-system/internal/controllers"
)

func CreateHTTPServer(cfg config.HttpConfig, tc *controllers.TransactionController) (*http.Server, error) {
	var (
		mainRouter = mux.NewRouter()
	)

	//Transaction handlers
	transactionHandlerFactory := NewTransactionHandlerFactory(tc)
	createTransactionHandler := securedHandler([]byte(cfg.JwtKey), handlers.ContentTypeHandler(transactionHandlerFactory.GetCreateHandler(), "application/json").ServeHTTP)
	mainRouter.HandleFunc(cfg.TransactionPath, createTransactionHandler).Methods(http.MethodPost)

	loggingHandler := handlers.LoggingHandler(os.Stdout, handlers.RecoveryHandler()(mainRouter))
	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           loggingHandler,
		ReadHeaderTimeout: cfg.ServerTimeout,
	}

	return srv, nil
}
