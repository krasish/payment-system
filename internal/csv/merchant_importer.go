package csv

import (
	"context"
	enc_csv "encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/krasish/payment-system/internal/common"

	"github.com/krasish/payment-system/internal/controllers"
)

type MerchantImporter struct {
	c *controllers.MerchantController
}

func NewMerchantImporter(c *controllers.MerchantController) *MerchantImporter {
	return &MerchantImporter{c: c}
}

func (i *MerchantImporter) Import(pathToCSVFile string) error {
	file, err := os.Open(filepath.Clean(pathToCSVFile))
	if err != nil {
		return fmt.Errorf("while opening file from %q: %w", pathToCSVFile, err)
	}
	defer common.CloseWithLogOnError(file)

	r := enc_csv.NewReader(file)
	dtos := make([]*controllers.Merchant, 0)
	for {
		record, err := r.Read()

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("while importing merchants: %w", err)
		}
		dto := new(controllers.Merchant)
		if err := dto.CSVUnmarshal(record); err != nil {
			return err
		}
		dtos = append(dtos, dto)
	}
	return i.c.CreateMerchants(context.Background(), dtos)
}
