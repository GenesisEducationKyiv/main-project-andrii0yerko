package core_test

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/model"
	"bitcoinrateapp/pkg/testenv"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
)

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

	from := "sender@test.org"
	password := ""
	host := "localhost"

	auth := core.NewAuthentication(from, password, host)
	client := core.NewSMTPClient(from, auth, host, smtpPort)
	formatter := core.NewPlainEmailFormatter(from)
	sender := core.NewEmailSender(client, formatter)

	rate := model.NewExchangeRate(1000, "coin", "currency")

	err = sender.SendRate("receiver@test.org", rate)
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
