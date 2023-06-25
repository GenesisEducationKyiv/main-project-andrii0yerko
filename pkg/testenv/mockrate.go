package testenv

import "context"

type MockRate struct {
	ExpectedRate float64
}

func (m *MockRate) Value(_ context.Context) (float64, error) {
	return m.ExpectedRate, nil
}

func (m *MockRate) Description() string {
	return "mock rate"
}
