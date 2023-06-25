package core_test

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/testenv"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
)

func readTotalMessages(url string) (int, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	var messages MessagesRepresentation
	err = json.NewDecoder(response.Body).Decode(&messages)
	if err != nil {
		return 0, err
	}

	return messages.Total, nil
}

type MessagesRepresentation struct {
	Total int `json:"total"`
}

func TestEmailSendIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	smtpPort, httpPort, err := testenv.SetupTestMailserver(t)
	if err != nil {
		t.Fatalf("Could not setup mailserver: %s", err)
	}

	host := "localhost"
	sender := core.NewEmailSender("sender@test.org", "", host, smtpPort)

	err = sender.Send("receiver@test.org", "Test", "Test")
	if err != nil {
		t.Fatalf("Could not send email: %s", err)
	}

	url := fmt.Sprintf("http://%s/api/v2/messages", net.JoinHostPort(host, httpPort))
	count, err := readTotalMessages(url)
	if err != nil {
		t.Fatalf("Could not read total messages: %s", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 message, got %d", count)
	}
}
