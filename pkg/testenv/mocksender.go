package testenv

type MockSender struct {
	ReceivedValues []string
	LastSubject    string
	LastMessage    string
}

func (m *MockSender) Send(receiver string, subject, message string) error {
	m.ReceivedValues = append(m.ReceivedValues, receiver)
	m.LastSubject = subject
	m.LastMessage = message
	return nil
}
