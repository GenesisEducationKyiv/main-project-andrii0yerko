package functional_test

import (
	"bitcoinrateapp/pkg/core"
	app "bitcoinrateapp/pkg/http"
	"bitcoinrateapp/pkg/testenv"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHTTPServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	smtpPort, _, err := testenv.SetupTestMailserver(t)
	if err != nil {
		t.Fatalf("Could not setup mailserver: %s", err)
	}

	receivers := testenv.NewTemporaryFileDB(t)
	sender := core.NewEmailSender("test@email.com", "", "localhost", smtpPort)
	rateRequester := &testenv.MockRate{ExpectedRate: 1000}

	service := core.NewService(receivers, rateRequester, sender)
	handler := app.NewExchangeRateHandler(service)
	addr := "localhost:3333"
	startServer(handler, addr, t)

	t.Run("get rate", func(t *testing.T) { testGetRate(addr, t) })
	t.Run("subscribe", func(t *testing.T) { testSubscribe(addr, t) })
	t.Run("duplicate subscribe", func(t *testing.T) { testDuplicateSubscribe(addr, t) })
	t.Run("send emails", func(t *testing.T) { testSendEmails(addr, t) })
}

func startServer(handler *app.ExchangeRateHandler, addr string, t *testing.T) {
	server := app.NewServer(handler, addr)

	// Start the server in a separate goroutine
	go func() {
		err := server.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Server error: %s", err)
		}
	}()
	time.Sleep(1 * time.Second)
	t.Cleanup(func() {
		err := server.Shutdown()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Shutdown error: %s", err)
		}
	})
}

func testGetRate(addr string, t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/rate", addr))
	if err != nil {
		t.Errorf("Request error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func testSubscribe(addr string, t *testing.T) {
	resp, err := http.PostForm(
		fmt.Sprintf("http://%s/subscribe", addr), map[string][]string{"email": {"test@test"}},
	)
	if err != nil {
		t.Errorf("Request error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func testDuplicateSubscribe(addr string, t *testing.T) {
	resp, err := http.PostForm(
		fmt.Sprintf("http://%s/subscribe", addr), map[string][]string{"email": {"test@test"}},
	)
	if err != nil {
		t.Errorf("Request error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

func testSendEmails(addr string, t *testing.T) {
	resp, err := http.Post(fmt.Sprintf("http://%s/sendEmails", addr), "application/json", nil)
	if err != nil {
		t.Errorf("Request error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
