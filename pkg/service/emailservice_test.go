package service_test

import (
	"bitcoinrateapp/pkg/model"
	"bitcoinrateapp/pkg/service"
	"bitcoinrateapp/pkg/testenv"
	"errors"
	"testing"
)

func TestServiceSubscribeSuccessfully(t *testing.T) {
	receiver := "abc@abc.test"
	subscriber, err := model.NewSubscriber(receiver)
	if err != nil {
		t.Fatal(err)
	}
	db := &testenv.MockDB{}
	btcservice := service.NewSenderService(db, nil)

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
	btcservice := service.NewSenderService(db, nil)

	err = btcservice.Subscribe(subscriber)
	if !errors.Is(err, expError) {
		t.Fatal(err)
	}
}

func TestServiceNotify(t *testing.T) {
	receivers := []string{"abc@abc.test", "abc2@abc.test"}
	rateValue := 100.0
	rate := model.NewExchangeRate(rateValue, "bitcoin", "uah")
	db := &testenv.MockDB{Memory: receivers}
	sender := &testenv.MockSender{}
	btcservice := service.NewSenderService(db, sender)

	err := btcservice.Notify(rate)

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
	if sender.LastRate.Value() != rateValue {
		t.Errorf("unexpected message: %s", sender.LastRate)
	}
}
