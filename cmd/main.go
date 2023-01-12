package main

import (
	"context"
	"log"
	"time"

	"github.com/krasish/payment-system/internal/csv"

	"github.com/krasish/payment-system/internal/views"

	"net/http"

	"github.com/krasish/payment-system/internal/config"
	"github.com/krasish/payment-system/internal/controllers"
	ps_http "github.com/krasish/payment-system/internal/http"
	"github.com/krasish/payment-system/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const ViewLayout = "bootstrap"

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("while reading configuration: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig.GetConnString()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	handleCSVImports(cfg, db)

	merchantStore := models.NewMerchantStore(db)
	merchantController := controllers.NewMerchantController(merchantStore)

	transactionStore := models.NewTransactionStore(db)
	transactionController := controllers.NewTransactionController(transactionStore, merchantStore)

	view, err := views.NewView(ViewLayout, cfg.ViewTemplatesPath)
	if err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	httpServer, err := ps_http.CreateHTTPServer(cfg.HttpConfig, transactionController, merchantController, view)
	if err != nil {
		log.Fatalln(err.Error())
	}
	transactionDeleter := transactionStore.GetPeriodicJobDeleter(time.Hour, cfg.DeletionJobInterval)
	go transactionDeleter(ctx)

	logrus.Infof("Running HTTP server on %s...", cfg.HttpConfig.Address)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Errorf("HTTP server ListenAndServe: %v", err)
	}

}

func handleCSVImports(cfg config.Config, db *gorm.DB) {
	if cfg.AdminsImportPath != "" {
		userStore := models.NewUserStore(db)
		userController := controllers.NewUserController(userStore)
		importer := csv.NewAdminImporter(userController)
		err := importer.Import(cfg.AdminsImportPath)
		if err != nil {
			log.Fatalf("while importing admins: %v", err.Error())
		}
	}
	if cfg.MerchantsImportPath != "" {
		merchantStore := models.NewMerchantStore(db)
		merchantController := controllers.NewMerchantController(merchantStore)
		importer := csv.NewMerchantImporter(merchantController)
		err := importer.Import(cfg.MerchantsImportPath)
		if err != nil {
			log.Fatalf("while importing merchants: %v", err.Error())
		}
	}
}
