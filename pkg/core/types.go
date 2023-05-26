package core

// An abstract storage which allows to read and add values
type Storage[T any] interface {
	Append(T) error
	GetRecords() ([]T, error)
}

// Abstract requester which allows to extract a specific value, and its description
type ValueRequester[T any] interface {
	GetValue() (T, error)
	GetDescription() string
}

// Defines behavior of sending data for the users
type Sender interface {
	Send(receiver string, subject string, message string) error
}
