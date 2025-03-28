package models_test

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/krasish/payment-system/internal/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const TransactionTestSchemaName = "payment_system_transaction_test"

var _ = Describe("Using NewTransction", func() {
	It("creates transaction without error for correct values", func() {
		t1, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "test@yahoo.com", "0888555885", 1, nil)
		Expect(err).To(BeNil())

		_, err = models.NewTransaction(uuid.Generate().String(), models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "test@yahoo.com", "0888555885", 1, &t1.ID)
		Expect(err).To(BeNil())
	})
	It("fails to create transaction when email is wrong", func() {
		_, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "", "0888555885", 1, nil)
		Expect(err).NotTo(BeNil())

		_, err = models.NewTransaction(uuid.Generate().String(), models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "www.google.com", "0888555885", 1, nil)
		Expect(err).NotTo(BeNil())
	})
	It("fails to create transaction when uuid is wrong", func() {
		_, err := models.NewTransaction("", models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "test@yahoo.com", "0888555885", 1, nil)
		Expect(err).NotTo(BeNil())

		_, err = models.NewTransaction("123", models.ToCurrency(500), models.TypeAuthorize, models.StatusApproved, "test@yahoo.com", "0888555885", 1, nil)
		Expect(err).NotTo(BeNil())
	})
})

var _ = Describe("Using TransactionStore", func() {
	var (
		merchantStore    *models.MerchantStore
		transactionStore *models.TransactionStore
		merchant         *models.Merchant
		err              error
		merchantEmail    = "mt@abv.bg"
		customerEmail    = "tc@mail.bg"
		customerPhone    = "0889787878"
	)

	BeforeEach(func() {
		_, err = sqlDB.Exec(fmt.Sprintf(SetSearchPathStatementFormat, TransactionTestSchemaName))
		Expect(err).To(BeNil())

		merchantStore = models.NewMerchantStore(gormDB)
		transactionStore = models.NewTransactionStore(gormDB)
		merchant, err = models.NewMerchant("Merchant With Transactions", "Hello!", merchantEmail, models.StatusActive)
		Expect(err).To(BeNil())

		err = merchantStore.CreateMerchant(context.Background(), merchant)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = merchantStore.DeleteMerchant(context.Background(), merchant.Email)
		Expect(err).To(BeNil())
	})

	Context("to create get and delete transactions", func() {
		It("works for Authorize -> Charge -> Refund flow", func() {
			transaction1, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeAuthorize, models.StatusApproved, customerEmail, customerPhone, merchant.UserID, nil)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction1)
			Expect(err).To(BeNil())

			transaction2, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeCharge, models.StatusApproved, customerEmail, customerPhone, merchant.UserID, &transaction1.ID)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction2)
			Expect(err).To(BeNil())

			returnedTransaction2, err := transactionStore.GetTransactionByUUID(context.Background(), transaction2.ExternalID)
			Expect(err).To(BeNil())
			compareTransactions([]*models.Transaction{transaction2}, []*models.Transaction{returnedTransaction2})

			transaction3, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeRefund, models.StatusApproved, customerEmail, customerPhone, merchant.UserID, &transaction2.ID)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction3)
			Expect(err).To(BeNil())

			transaction2.Status = models.StatusRefunded
			transactions, err := transactionStore.GetAllTransactions(context.Background())
			Expect(err).To(BeNil())

			compareTransactions([]*models.Transaction{transaction1, transaction2, transaction3}, transactions)

			//Try to delete all of them to check for issues with db cascade
			err = transactionStore.DeleteTransaction(context.Background(), transaction1)
			Expect(err).To(BeNil())

			err = transactionStore.DeleteTransaction(context.Background(), transaction2)
			Expect(err).To(BeNil())

			err = transactionStore.DeleteTransaction(context.Background(), transaction3)
			Expect(err).To(BeNil())
		})

		It("works for Authorize -> Reverse flow", func() {
			transaction1, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeAuthorize, models.StatusApproved, customerEmail, customerPhone, merchant.UserID, nil)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction1)
			Expect(err).To(BeNil())

			transaction2, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeReversal, models.StatusReversed, customerEmail, customerPhone, merchant.UserID, &transaction1.ID)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction2)
			Expect(err).To(BeNil())

			transaction1.Status = models.StatusReversed

			transactions, err := transactionStore.GetAllTransactions(context.Background())
			Expect(err).To(BeNil())
			compareTransactions([]*models.Transaction{transaction1, transaction2}, transactions)

			err = transactionStore.DeleteTransaction(context.Background(), transaction1)
			Expect(err).To(BeNil())

			err = transactionStore.DeleteTransaction(context.Background(), transaction2)
			Expect(err).To(BeNil())
		})

		It("does not create transactions for inactive merchants", func() {
			inactiveMerchant, err := models.NewMerchant("Merchant With Transactions", "Hello!", "inactive_merchant@gmail.com", models.StatusInactive)
			Expect(err).To(BeNil())

			err = merchantStore.CreateMerchant(context.Background(), inactiveMerchant)
			Expect(err).To(BeNil())

			transaction1, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeAuthorize, models.StatusApproved, customerEmail, customerPhone, inactiveMerchant.UserID, nil)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction1)
			Expect(err).NotTo(BeNil())
		})
	})

	Context("for periodic transactions deletion", func() {
		It("deletes transactions", func() {
			transaction1, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeAuthorize, models.StatusApproved, customerEmail, customerPhone, merchant.UserID, nil)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction1)
			Expect(err).To(BeNil())

			transaction2, err := models.NewTransaction(uuid.Generate().String(), models.ToCurrency(800), models.TypeCharge, models.StatusApproved, customerEmail, customerPhone, merchant.UserID, &transaction1.ID)
			Expect(err).To(BeNil())

			err = transactionStore.CreateTransaction(context.Background(), transaction2)
			Expect(err).To(BeNil())

			transactionDeleter := transactionStore.GetPeriodicJobDeleter(100*time.Millisecond, 100*time.Millisecond)
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			DeferCleanup(cancel)

			go transactionDeleter(ctx)
			<-ctx.Done()

			//TODO: There is an issue causing the gormDB to use another schema. Fix it
			_, err = sqlDB.Exec(fmt.Sprintf(SetSearchPathStatementFormat, TransactionTestSchemaName))
			Expect(err).To(BeNil())

			transactions, err := transactionStore.GetAllTransactions(context.Background())
			Expect(err).To(BeNil())
			Expect(transactions).To(BeEmpty())
		})
	})
})

func compareTransactions(expected []*models.Transaction, actual []*models.Transaction) {
	type Validation struct {
		ID            uint
		CustomerMail  string
		CustomerPhone string
		Status        models.TransactionStatus
		Type          models.TransactionType
		MerchantID    uint
		BelongsToID   *uint
	}
	Expect(expected).To(HaveLen(len(actual)))
	expectV, actualV := make([]Validation, 0, len(expected)), make([]Validation, 0, len(expected))

	for i := 0; i < len(expected); i++ {
		//Works since they have the same length
		expectV, actualV = append(expectV, Validation{
			ID:            expected[i].ID,
			Status:        expected[i].Status,
			CustomerMail:  expected[i].CustomerEmail,
			CustomerPhone: expected[i].CustomerPhone,
			Type:          expected[i].Type,
			MerchantID:    expected[i].MerchantID,
			BelongsToID:   expected[i].BelongsToID,
		}), append(actualV, Validation{
			ID:            actual[i].ID,
			Status:        actual[i].Status,
			Type:          actual[i].Type,
			CustomerMail:  actual[i].CustomerEmail,
			CustomerPhone: actual[i].CustomerPhone,
			MerchantID:    actual[i].MerchantID,
			BelongsToID:   actual[i].BelongsToID,
		})
	}
	Expect(expectV).To(ContainElements(actualV))
}
