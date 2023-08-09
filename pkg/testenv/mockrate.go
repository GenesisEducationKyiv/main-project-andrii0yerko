package testenv

import (
	"bitcoinrateapp/pkg/model"
	"context"
)

type MockRate struct {
	ExpectedRate float64
}

func (m *MockRate) Value(_ context.Context, _, _ string) (model.Rate, error) {
	return model.NewExchangeRate(m.ExpectedRate, "coin", "currency"), nil
}
