package utils

type numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Min[T numeric](a T, b ...T) T {
	for _, x := range b {
		if a > x {
			a = x
		}
	}
	return a
}

func Max[T numeric](a T, b ...T) T {
	for _, x := range b {
		if a < x {
			a = x
		}
	}
	return a
}
