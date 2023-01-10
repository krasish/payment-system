package csv

import (
	"context"
	enc_csv "encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/krasish/payment-system/internal/controllers"
)

type AdminImporter struct {
	c *controllers.UserController
}

func NewAdminImporter(c *controllers.UserController) *AdminImporter {
	return &AdminImporter{c: c}
}

func (i *AdminImporter) Import(pathToCSVFile string) error {
	file, err := os.Open(pathToCSVFile)
	if err != nil {
		return fmt.Errorf("while opening file from %q: %w", pathToCSVFile, err)
	}
	defer file.Close()

	r := enc_csv.NewReader(file)
	dtos := make([]*controllers.User, 0)
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("while importing merchants: %v", err)
		}
		dto := new(controllers.User)
		if err := dto.CSVUnmarshal(record); err != nil {
			return err
		}
		dtos = append(dtos, dto)
	}
	return i.c.CreateUsers(context.Background(), dtos)
}
