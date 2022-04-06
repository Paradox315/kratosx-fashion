package math

type Number interface {
	int64 | int32 | int16 | int8 | int | uint64 | uint32 | uint16 | uint8 | uint | float64 | float32 | ~string
}

func Max[T Number](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Min[T Number](x, y T) T {
	if x < y {
		return x
	}
	return y
}
