package testenv

import "context"

type MockRate struct {
	ExpectedRate float64
}

func (m *MockRate) Value(_ context.Context, _, _ string) (float64, error) {
	return m.ExpectedRate, nil
}
