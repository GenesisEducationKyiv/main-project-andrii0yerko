package testenv

type MockErrorDB struct {
	ExpectedError error
}

func (m *MockErrorDB) Records() ([]string, error) {
	return nil, nil
}

func (m *MockErrorDB) Append(_ string) error {
	return m.ExpectedError
}
