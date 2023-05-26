package core

type Storage[T any] interface {
	Append(T) error
	GetRecords() ([]T, error)
}

type ValueRequester interface {
	GetValue() (float64, error)
}
