package testenv

import (
	"bitcoinrateapp/pkg/model"
)

type MockSender struct {
	ReceivedValues []string
	LastRate       model.Rate
}

func (m *MockSender) SendRate(receiver string, rate model.Rate) error {
	m.ReceivedValues = append(m.ReceivedValues, receiver)
	m.LastRate = rate
	return nil
}
