package models_test

import (
	"context"
	"fmt"

	"github.com/krasish/payment-system/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	UserTestSchemaName           = "payment_system_user_test"
	SetSearchPathStatementFormat = "SET SEARCH_PATH TO %q; "
)

var (
	adminUser     = models.NewUser(models.RoleAdmin, models.StatusActive)
	merchantUser  = models.NewUser(models.RoleMerchant, models.StatusActive)
	inactiveUsers = []*models.User{
		models.NewUser(models.RoleAdmin, models.StatusInactive),
		models.NewUser(models.RoleMerchant, models.StatusInactive),
	}
)

var _ = Describe("Using UserStore", func() {
	var (
		userStore *models.UserStore
		err       error
	)

	BeforeEach(func() {
		_, err = sqlDB.Exec(fmt.Sprintf(SetSearchPathStatementFormat, UserTestSchemaName))
		Expect(err).To(BeNil())

		userStore = models.NewUserStore(gormDB)
	})

	//Notice the Serial decorator
	Context("to create and get users works", Serial, func() {
		It("creates admin user successfully", func() {
			err = userStore.CreateUser(context.Background(), adminUser)
			Expect(err).To(BeNil())
		})

		It("creates merchant user successfully", func() {
			err = userStore.CreateUser(context.Background(), merchantUser)
			Expect(err).To(BeNil())
		})

		It("creates list of inactive users successfully", func() {
			err = userStore.CreateUsers(context.Background(), inactiveUsers)
			Expect(err).To(BeNil())
		})

		It("gets all previously created users successfully", func() {
			users, err := userStore.GetAllUsers(context.Background())
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(4))
			Expect(users).Should(ContainElements(BeComparableTo(adminUser), BeComparableTo(merchantUser)))
			Expect(users).Should(ContainElements(BeComparableTo(inactiveUsers[0]), BeComparableTo(inactiveUsers[1])))
		})

	})

})
