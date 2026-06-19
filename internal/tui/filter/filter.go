package filter

type FilterAble[T any] interface {
	Matches(query string) bool
	GetItem() T
}
