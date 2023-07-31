package testenv

import "bitcoinrateapp/pkg/rateclient"

type MockSender struct {
	ReceivedValues []string
	LastRate       rateclient.Rate
}

func (m *MockSender) SendRate(receiver string, rate rateclient.Rate) error {
	m.ReceivedValues = append(m.ReceivedValues, receiver)
	m.LastRate = rate
	return nil
}
