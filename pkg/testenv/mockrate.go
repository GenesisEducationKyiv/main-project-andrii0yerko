package testenv

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/model"
	"context"
)

type MockRate struct {
	ExpectedRate float64
}

func (m *MockRate) Value(_ context.Context, _, _ string) (core.Rate, error) {
	return model.NewExchangeRate(m.ExpectedRate, "coin", "currency"), nil
}
