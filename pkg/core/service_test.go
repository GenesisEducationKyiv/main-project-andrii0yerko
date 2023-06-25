package core_test

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/testenv"
	"errors"
	"fmt"
	"testing"
)

type MockDB struct {
	records []string
}

func (m *MockDB) Records() ([]string, error) {
	return m.records, nil
}

func (m *MockDB) Append(value string) error {
	m.records = append(m.records, value)
	return nil
}

type MockFilledDB struct {
	expectedError error
}

func (m *MockFilledDB) Records() ([]string, error) {
	return nil, nil
}

func (m *MockFilledDB) Append(_ string) error {
	return m.expectedError
}

type MockSender struct {
	receivedValues []string
	lastSubject    string
	lastMessage    string
}

func (m *MockSender) Send(receiver string, subject, message string) error {
	m.receivedValues = append(m.receivedValues, receiver)
	m.lastSubject = subject
	m.lastMessage = message
	return nil
}

func TestServiceRate(t *testing.T) {
	rate := 100.0
	service := core.NewService(nil, &testenv.MockRate{ExpectedRate: rate}, nil)
	actualRate, err := service.ExchangeRate()
	if err != nil {
		t.Fatal(err)
	}
	if actualRate != rate {
		t.Errorf("unexpected rate: %f", actualRate)
	}
}

func TestServiceSubscribeSuccessfully(t *testing.T) {
	receiver := "abc@abc.test"
	db := &MockDB{}
	service := core.NewService(db, nil, nil)

	err := service.Subscribe(receiver)

	if err != nil {
		t.Fatal(err)
	}
	if len(db.records) != 1 {
		t.Errorf("unexpected records count: %d", len(db.records))
	}
	if db.records[0] != receiver {
		t.Errorf("unexpected record: %s", db.records[0])
	}
}

func TestServiceSubscribeError(t *testing.T) {
	receiver := "abc@abc.test"
	expError := core.ErrIsDuplicate
	db := &MockFilledDB{expectedError: expError}
	service := core.NewService(db, nil, nil)

	err := service.Subscribe(receiver)

	if !errors.Is(err, expError) {
		t.Fatal(err)
	}
}

func TestServiceNotify(t *testing.T) {
	// receiver := "abc@abc.test"
	receivers := []string{"abc@abc.test", "abc2@abc.test"}
	rate := 100.0
	db := &MockDB{records: receivers}
	rateProvider := &testenv.MockRate{ExpectedRate: rate}
	sender := &MockSender{}
	service := core.NewService(db, rateProvider, sender)

	err := service.Notify()

	if err != nil {
		t.Fatal(err)
	}
	if len(sender.receivedValues) != len(receivers) {
		t.Errorf("unexpected receivers count: %d", len(sender.receivedValues))
	}

	if sender.receivedValues[0] != receivers[0] {
		t.Errorf("unexpected receiver: %s", sender.receivedValues[0])
	}
	if sender.receivedValues[1] != receivers[1] {
		t.Errorf("unexpected receiver: %s", sender.receivedValues[1])
	}
	if sender.lastSubject != rateProvider.Description() {
		t.Errorf("unexpected subject: %s", sender.lastSubject)
	}
	if sender.lastMessage != fmt.Sprintf("%f", rate) {
		t.Errorf("unexpected message: %s", sender.lastMessage)
	}
}
