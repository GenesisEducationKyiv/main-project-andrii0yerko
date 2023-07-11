package service_test

import (
	"bitcoinrateapp/pkg/model"
	"bitcoinrateapp/pkg/service"
	"bitcoinrateapp/pkg/testenv"
	"errors"
	"testing"
)

func TestServiceRate(t *testing.T) {
	rate := 100.0
	btcservice := service.NewService(nil, &testenv.MockRate{ExpectedRate: rate}, nil, "bitcoin", "uah")
	actualRate, err := btcservice.ExchangeRate()
	if err != nil {
		t.Fatal(err)
	}
	if actualRate != rate {
		t.Errorf("unexpected rate: %f", actualRate)
	}
}

func TestServiceSubscribeSuccessfully(t *testing.T) {
	receiver := "abc@abc.test"
	subscriber, err := model.NewSubscriber(receiver)
	if err != nil {
		t.Fatal(err)
	}
	db := &testenv.MockDB{}
	btcservice := service.NewService(db, nil, nil, "bitcoin", "uah")

	err = btcservice.Subscribe(subscriber)
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
	subscriber, err := model.NewSubscriber(receiver)
	if err != nil {
		t.Fatal(err)
	}
	expError := service.ErrIsDuplicate
	db := &testenv.MockErrorDB{ExpectedError: expError}
	btcservice := service.NewService(db, nil, nil, "bitcoin", "uah")

	err = btcservice.Subscribe(subscriber)
	if !errors.Is(err, expError) {
		t.Fatal(err)
	}
}

func TestServiceNotify(t *testing.T) {
	receivers := []string{"abc@abc.test", "abc2@abc.test"}
	rate := 100.0
	db := &testenv.MockDB{Memory: receivers}
	rateProvider := &testenv.MockRate{ExpectedRate: rate}
	sender := &testenv.MockSender{}
	btcservice := service.NewService(db, rateProvider, sender, "bitcoin", "uah")

	err := btcservice.Notify()

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
	if sender.LastRate.Value() != rate {
		t.Errorf("unexpected message: %s", sender.LastRate)
	}
}
