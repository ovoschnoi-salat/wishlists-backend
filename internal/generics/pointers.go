package generics

func GetPtr[T any](t T) *T {
	return &t
}

func FromPtr[T any](t *T) T {
	if t == nil {
		var zeroValue T
		return zeroValue
	}
	return *t
}
