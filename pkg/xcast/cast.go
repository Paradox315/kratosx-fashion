package xcast

import "github.com/spf13/cast"

func ToUint64Slice[T any](arr []T) (uint64s []uint64) {
	uint64s = make([]uint64, len(arr))
	for i := range arr {
		uint64s[i] = cast.ToUint64(arr[i])
	}
	return
}

func ToUintSlice[T any](arr []T) (uints []uint) {
	uints = make([]uint, len(arr))
	for i := range arr {
		uints[i] = cast.ToUint(arr[i])
	}
	return
}
