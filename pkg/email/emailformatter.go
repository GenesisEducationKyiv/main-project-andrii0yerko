package email

import (
	"bitcoinrateapp/pkg/model"
	"fmt"
	"strings"
)

type PlainEmailFormatter struct {
	from string
}

func NewPlainEmailFormatter(from string) PlainEmailFormatter {
	return PlainEmailFormatter{from: from}
}

func (e PlainEmailFormatter) Format(receiver string, rate model.Rate) string {
	subject := e.getDescription(rate.Coin(), rate.Currency())
	text := fmt.Sprintf("%f", rate.Value())
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", e.from, receiver, subject, text)
	return message
}

func (e PlainEmailFormatter) getDescription(coin, currency string) string {
	cointTitle := strings.ToUpper(coin[:1]) + strings.ToLower(coin[1:])
	return fmt.Sprintf("%s to %s Exchange Rate", cointTitle, strings.ToUpper(currency))
}
