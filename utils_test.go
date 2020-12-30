package hego

import "testing"

func TestConvertBool(t *testing.T) {
	arr := make([]interface{}, 0)
	for i := 0; i < 3; i++ {
		arr = append(arr, true)
	}
	ConvertBool(arr)
}

func TestConvertFloat64(t *testing.T) {
	arr := make([]interface{}, 0)
	for i := 0; i < 3; i++ {
		arr = append(arr, 1.0)
	}
	ConvertFloat64(arr)
}

func TestConvertInt(t *testing.T) {
	arr := make([]interface{}, 0)
	for i := 0; i < 3; i++ {
		arr = append(arr, 1)
	}
	ConvertInt(arr)
}
