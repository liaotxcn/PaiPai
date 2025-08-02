package utils

func ConvertToInt8[T int8 | int32 | int64 | int](value T) *int8 {
	v := int8(value)
	return &v
}
