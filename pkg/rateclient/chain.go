package rateclient

import (
	"bitcoinrateapp/pkg/model"
	"context"
)

type requester interface {
	Value(ctx context.Context, coin, currency string) (model.Rate, error)
}

type RequesterChain struct {
	element requester
	next    requester
}

func NewRequesterChain(element requester) *RequesterChain {
	return &RequesterChain{
		element: element,
	}
}

func (c *RequesterChain) Value(ctx context.Context, coin, currency string) (model.Rate, error) {
	rate, err := c.element.Value(ctx, coin, currency)
	if err != nil {
		if c.next != nil {
			return c.next.Value(ctx, coin, currency)
		}
		return nil, err
	}
	return rate, nil
}

func (c *RequesterChain) SetNext(next requester) {
	c.next = next
}
