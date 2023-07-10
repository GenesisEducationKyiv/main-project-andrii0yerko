package model

type Rate interface {
	Value() float64
	Coin() string
	Currency() string
}
