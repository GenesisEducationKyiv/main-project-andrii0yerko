package testenv

type MockDB struct {
	Memory []string
}

func (m *MockDB) Records() ([]string, error) {
	return m.Memory, nil
}

func (m *MockDB) Append(value string) error {
	m.Memory = append(m.Memory, value)
	return nil
}
