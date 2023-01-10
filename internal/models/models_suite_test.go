package models_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm/schema"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"github.com/krasish/payment-system/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
)

var (
	pathToDDL          = "/../../sql/schema.sql"
	postgresImage      = "postgres:15"
	testDatabaseConfig = config.DatabaseConfig{
		User:     "test-user",
		Password: "test-pas$word",
		Host:     "localhost",
		Port:     "5432",
		Name:     "payment_system",
		SSLMode:  "disable",
	}
	postgresContainer testcontainers.Container
	sqlDB             *sql.DB
	gormDB            *gorm.DB
	testDurationLimit = time.Minute
	schemaNames       = []string{UserTestSchemaName, MerchantTestSchemaName, TransactionTestSchemaName}
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	var (
		tcpSuffix   = "/tcp"
		exposedPort = testDatabaseConfig.Port + tcpSuffix
		err         error
		ctx, cancel = context.WithTimeout(context.Background(), testDurationLimit)
	)

	workDir, err := os.Getwd()
	Expect(err).To(BeNil())

	schemaSQL, err := os.ReadFile(workDir + pathToDDL)
	Expect(err).To(BeNil())

	req := testcontainers.ContainerRequest{
		Name:         "metadata-directory-tests-postgres",
		User:         "postgres",
		Image:        postgresImage,
		ExposedPorts: []string{exposedPort},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     testDatabaseConfig.User,
			"POSTGRES_PASSWORD": testDatabaseConfig.Password,
			"POSTGRES_DB":       testDatabaseConfig.Name,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(exposedPort)),
	}
	postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).To(BeNil())
	DeferCleanup(func() {
		err = postgresContainer.Terminate(ctx)
		Expect(err).To(BeNil())
		cancel()
	})
	Expect(postgresContainer.IsRunning()).To(BeTrue())

	port, err := postgresContainer.MappedPort(ctx, nat.Port(testDatabaseConfig.Port))
	Expect(err).To(BeNil())
	testDatabaseConfig.Port = strings.TrimSuffix(string(port), tcpSuffix)

	sqlDB, err = sql.Open("postgres", testDatabaseConfig.GetConnString())
	Expect(err).To(BeNil())

	err = sqlDB.Ping()
	Expect(err).To(BeNil())

	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	Expect(err).To(BeNil())

	for _, schemaName := range schemaNames {
		migration := addSchemaToMigration(string(schemaSQL), schemaName)
		_, err = sqlDB.Exec(migration)
		Expect(err).To(BeNil())
	}

	zone, _ := time.Now().Zone()
	_, err = sqlDB.Exec(fmt.Sprintf("SET TIME ZONE %q;", zone))
	Expect(err).To(BeNil())
})

func addSchemaToMigration(migration, schema string) string {
	prefix := fmt.Sprintf("BEGIN TRANSACTION;\n CREATE SCHEMA %s;\n SET search_path TO %s;\n ;COMMIT;\n", schema, schema)
	return prefix + migration
}
