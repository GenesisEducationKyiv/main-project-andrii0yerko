package model

import "errors"

var ErrEmptyEmail = errors.New("empty email")

type Subscriber struct {
	email string
}

func NewSubscriber(email string) (*Subscriber, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}
	return &Subscriber{email: email}, nil
}

func (u *Subscriber) Email() string {
	return u.email
}
