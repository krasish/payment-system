package models_test

import (
	"github.com/krasish/payment-system/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Using Currency values", func() {
	Context("from float64", func() {
		It("works with whole values as expected", func() {
			Expect(models.ToCurrency(0.0)).To(BeEquivalentTo(0))
			Expect(models.ToCurrency(5.0)).To(BeEquivalentTo(500))
		})
		It("works with whole values with tens after decimal point as expected", func() {
			Expect(models.ToCurrency(0.1)).To(BeEquivalentTo(10))
			Expect(models.ToCurrency(5.1)).To(BeEquivalentTo(510))
		})
		It("works with whole values with tens and ones after decimal point as expected", func() {
			Expect(models.ToCurrency(0.01)).To(BeEquivalentTo(1))
			Expect(models.ToCurrency(0.16)).To(BeEquivalentTo(16))
			Expect(models.ToCurrency(5.12)).To(BeEquivalentTo(512))
		})
		It("works correctly with more precise values", func() {
			Expect(models.ToCurrency(5.123)).To(BeEquivalentTo(512))
			Expect(models.ToCurrency(5.126)).To(BeEquivalentTo(513))
		})
	})

	Context("to float64", func() {
		It("works with whole values as expected", func() {
			Expect(models.Currency(0).Float64()).Should(Equal(0.0))
			Expect(models.Currency(500).Float64()).Should(Equal(5.0))
		})
		It("works with whole values with tens after decimal point as expected", func() {
			Expect(models.Currency(10).Float64()).Should(Equal(0.1))
			Expect(models.Currency(510).Float64()).Should(Equal(5.1))
		})
		It("works with whole values with tens and ones after decimal point as expected", func() {
			Expect(models.Currency(1).Float64()).Should(Equal(0.01))
			Expect(models.Currency(16).Float64()).Should(Equal(0.16))
			Expect(models.Currency(512).Float64()).Should(Equal(5.12))
		})
	})

	Context("in a round trip between types", func() {
		It("survives", func() {
			floats := []float64{0.0, 0.1, 0.12, 5.0, 5.1, 5.12, 5.10, 5.01, 10.2, 10.1, 10.23}
			for _, f := range floats {
				curr := models.ToCurrency(f)
				Expect(curr.Float64()).To(Equal(f))
			}
		})
	})
})
