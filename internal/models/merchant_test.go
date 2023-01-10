package models_test

import (
	"context"
	"fmt"

	"github.com/docker/distribution/uuid"

	"github.com/krasish/payment-system/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const MerchantTestSchemaName = "payment_system_merchant_test"

var _ = Describe("Using NewMerchant", func() {
	It("creates merchant without error for correct values", func() {
		_, err := models.NewMerchant("Sample Merchant", "Desctiption", "test@yahoo.com", models.StatusActive)
		Expect(err).To(BeNil())
	})
	It("fails to create merchant when email is wrong", func() {
		_, err := models.NewMerchant("Sample Merchant", "Desctiption", "", models.StatusActive)
		Expect(err).NotTo(BeNil())
		_, err = models.NewMerchant("Sample Merchant", "Desctiption", "www.google.com", models.StatusActive)
		Expect(err).NotTo(BeNil())
	})
})

var (
	merchant, _  = models.NewMerchant("Merchant One", "Merchant One Description", "merchant1@abv.bg", models.StatusActive)
	merchant2, _ = models.NewMerchant("Merchant Two", "Merchant Two Description", "merchant2@abv.bg", models.StatusActive)
	merchant3, _ = models.NewMerchant("Merchant Three", "Merchant Three Description", "merchant3@abv.bg", models.StatusInactive)

	merchantWithTransactions, _ = models.NewMerchant("Merchant With Transaction", "Merchant With Transaction Description", "merchant4@abv.bg", models.StatusActive)
)

var _ = Describe("Using MerchantStore", func() {
	var (
		merchantStore    *models.MerchantStore
		transactionStore *models.TransactionStore
		err              error
	)

	BeforeEach(func() {
		_, err = sqlDB.Exec(fmt.Sprintf(SetSearchPathStatementFormat, MerchantTestSchemaName))
		Expect(err).To(BeNil())

		merchantStore = models.NewMerchantStore(gormDB)
		transactionStore = models.NewTransactionStore(gormDB)
	})

	//Notice the Serial decorator
	Context("to create and get merchants", Serial, func() {
		It("creates proper merchant successfully", func() {
			err := merchantStore.CreateMerchant(context.Background(), merchant)
			Expect(err).To(BeNil())
		})
		It("creates list of merchants successfully", func() {
			err := merchantStore.CreateMerchants(context.Background(), []*models.Merchant{merchant2, merchant3})
			Expect(err).To(BeNil())
		})
		It("gets all previously created merchants", func() {
			merchants, err := merchantStore.GetAllMerchants(context.Background())
			Expect(err).To(BeNil())
			merchant.Transactions = make([]models.Transaction, 0)
			merchant2.Transactions = make([]models.Transaction, 0)
			merchant3.Transactions = make([]models.Transaction, 0)
			Expect(merchants).Should(ContainElements(BeComparableTo(merchant), BeComparableTo(merchant2), BeComparableTo(merchant3)))
		})
	})

	Context("to create & get merchants with transactions", Serial, func() {
		var (
			transaction1 *models.Transaction
			transaction2 *models.Transaction
		)
		It("creates merchant and its transactions successfully", func() {
			err := merchantStore.CreateMerchant(context.Background(), merchantWithTransactions)
			Expect(err).To(BeNil())

			transaction1, _ = models.NewTransaction(uuid.Generate().String(), models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "cusotmer1@yahoo.com", "0889998989", merchantWithTransactions.UserID, nil)
			err = transactionStore.CreateTransaction(context.Background(), transaction1)
			Expect(err).To(BeNil())

			transaction2, _ = models.NewTransaction(uuid.Generate().String(), models.ToCurrency(600), models.TypeCharge, models.StatusApproved, "cusotmer1@yahoo.com", "0889998989", merchantWithTransactions.UserID, &transaction1.ID)
			err = transactionStore.CreateTransaction(context.Background(), transaction2)
			Expect(err).To(BeNil())
		})

		It("retrieves merchant by id with its transactions successfully", func() {
			res, err := merchantStore.GetMerchantById(context.Background(), merchantWithTransactions.UserID)
			Expect(err).To(BeNil())
			Expect(res.UserID).To(Equal(merchantWithTransactions.UserID))
			Expect(res.Name).To(Equal(merchantWithTransactions.Name))
			Expect(res.Description).To(Equal(merchantWithTransactions.Description))

			returnedTransactions := make([]*models.Transaction, 0)
			for i := range res.Transactions {
				returnedTransactions = append(returnedTransactions, &res.Transactions[i])
			}
			compareTransactions([]*models.Transaction{transaction1, transaction2}, returnedTransactions)

		})

	})

})
