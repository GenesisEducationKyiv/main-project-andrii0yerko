package core_test

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/testenv"
	"errors"
	"fmt"
	"testing"
)

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
	db := &testenv.MockDB{}
	service := core.NewService(db, nil, nil)

	err := service.Subscribe(receiver)

	if err != nil {
		t.Fatal(err)
	}
	if len(db.Memory) != 1 {
		t.Errorf("unexpected records count: %d", len(db.Memory))
	}
	if db.Memory[0] != receiver {
		t.Errorf("unexpected record: %s", db.Memory[0])
	}
}

func TestServiceSubscribeError(t *testing.T) {
	receiver := "abc@abc.test"
	expError := core.ErrIsDuplicate
	db := &testenv.MockErrorDB{ExpectedError: expError}
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
	db := &testenv.MockDB{Memory: receivers}
	rateProvider := &testenv.MockRate{ExpectedRate: rate}
	sender := &testenv.MockSender{}
	service := core.NewService(db, rateProvider, sender)

	err := service.Notify()

	if err != nil {
		t.Fatal(err)
	}
	if len(sender.ReceivedValues) != len(receivers) {
		t.Errorf("unexpected receivers count: %d", len(sender.ReceivedValues))
	}

	if sender.ReceivedValues[0] != receivers[0] {
		t.Errorf("unexpected receiver: %s", sender.ReceivedValues[0])
	}
	if sender.ReceivedValues[1] != receivers[1] {
		t.Errorf("unexpected receiver: %s", sender.ReceivedValues[1])
	}
	if sender.LastSubject != rateProvider.Description() {
		t.Errorf("unexpected subject: %s", sender.LastSubject)
	}
	if sender.LastMessage != fmt.Sprintf("%f", rate) {
		t.Errorf("unexpected message: %s", sender.LastMessage)
	}
}
