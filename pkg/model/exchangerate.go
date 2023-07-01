package model

type ExchangeRate struct {
	value    float64
	coin     string
	currency string
}

func NewExchangeRate(value float64, coin, currency string) *ExchangeRate {
	return &ExchangeRate{
		value:    value,
		coin:     coin,
		currency: currency,
	}
}

func (r *ExchangeRate) Value() float64 {
	return r.value
}

func (r *ExchangeRate) Coin() string {
	return r.coin
}

func (r *ExchangeRate) Currency() string {
	return r.currency
}
