package core

type Storage[T any] interface {
	Append(T) error
	GetRecords() ([]T, error)
}

type ValueRequester[T any] interface {
	GetValue() (T, error)
}

type Sender interface {
	Send(string, string, string) error
}
