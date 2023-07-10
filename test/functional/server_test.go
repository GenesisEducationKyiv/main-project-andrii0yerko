package functional_test

import (
	"bitcoinrateapp/pkg/app"
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/testenv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	receivers, file := testenv.NewTemporaryFileDB(t)

	from := "test@email.com"
	password := ""
	host := "localhost"

	auth := core.NewAuthentication(from, password, host)
	client := core.NewSMTPClient(from, auth, host, smtpPort)
	formatter := core.NewPlainEmailFormatter(from)
	sender := core.NewEmailSender(client, formatter)
	rateRequester := &testenv.MockRate{ExpectedRate: 1000}

	service := core.NewService(receivers, rateRequester, sender)
	handler := app.NewExchangeRateHandler(service)
	addr := "localhost:3333"
	startServer(handler, addr, t)

	t.Run("get rate", func(t *testing.T) { testGetRate(addr, t) })
	t.Run("subscribe", func(t *testing.T) {
		runWithTransaction(file, func() { testSubscribe(addr, t) })
	})
	t.Run("duplicate subscribe", func(t *testing.T) {
		runWithTransaction(file, func() { testDuplicateSubscribe(addr, receivers, t) })
	})
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Read body error: %s", err)
	}
	if len(body) == 0 {
		t.Errorf("Empty body")
	}
	for _, b := range body {
		if b != '.' && (b < '0' || b > '9') {
			t.Errorf("Body contains non-digit: %c", b)
		}
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

func testDuplicateSubscribe(addr string, db core.Storage[string], t *testing.T) {
	email := "test@test"
	err := db.Append(email)
	if err != nil {
		t.Errorf("Setup DB data: %s", err)
	}
	resp, err := http.PostForm(
		fmt.Sprintf("http://%s/subscribe", addr), map[string][]string{"email": {email}},
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

func runWithTransaction(file *os.File, fn func()) {
	currentState, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read file state:", err)
	}

	fn()

	err = file.Truncate(0)
	if err != nil {
		log.Fatal("Failed to truncate file:", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal("Failed to seek file:", err)
	}
	_, err = file.Write(currentState)
	if err != nil {
		log.Fatal("Failed to write file state:", err)
	}
}
