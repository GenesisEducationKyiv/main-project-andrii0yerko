package testenv

import (
	"bitcoinrateapp/pkg/model"
	"bitcoinrateapp/pkg/rateclient"
	"context"
)

type MockRate struct {
	ExpectedRate float64
}

func (m *MockRate) Value(_ context.Context, _, _ string) (rateclient.Rate, error) {
	return model.NewExchangeRate(m.ExpectedRate, "coin", "currency"), nil
}
