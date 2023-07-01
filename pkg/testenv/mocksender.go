package testenv

import "bitcoinrateapp/pkg/core"

type MockSender struct {
	ReceivedValues []string
	LastRate       core.Rate
}

func (m *MockSender) SendRate(receiver string, rate core.Rate) error {
	m.ReceivedValues = append(m.ReceivedValues, receiver)
	m.LastRate = rate
	return nil
}
