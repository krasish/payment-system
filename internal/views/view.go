package views

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/krasish/payment-system/internal/controllers"
)

type Merchant struct {
	*controllers.Merchant
	Transactions []*controllers.Transaction
}

func NewMerchant(merchant *controllers.Merchant) *Merchant {
	return &Merchant{Merchant: merchant, Transactions: make([]*controllers.Transaction, 0)}
}

type MerchantsData map[string]*Merchant

func NewMerchantsData(merchants []*controllers.Merchant, transactions []*controllers.Transaction) MerchantsData {
	md := make(MerchantsData, len(merchants))
	for i, merchant := range merchants {
		md[merchant.Email] = NewMerchant(merchants[i])
	}
	for i, transaction := range transactions {
		md[transaction.MerchantEmail].Transactions = append(md[transaction.MerchantEmail].Transactions, transactions[i])
	}
	return md
}

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(layout, templatesDir string) (*View, error) {
	files, err := filepath.Glob(templatesDir + "/*.gohtml")
	if err != nil {
		logrus.Errorf("while opening files with templates: %v", err)
		return nil, fmt.Errorf("while opening files with templates: %w", err)
	}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		logrus.Errorf("while parsing file tempaltes: %v", err)
		return nil, fmt.Errorf("while parsing file tempaltes: %w", err)
	}
	return &View{Template: tpl, Layout: layout}, nil
}

func (v *View) Render(w http.ResponseWriter, merchants MerchantsData) error {
	return v.Template.ExecuteTemplate(w, v.Layout, merchants)
}
