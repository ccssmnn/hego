package hego

// ConvertBool uses type assertion to convert an interface slice to a bool slice
func ConvertBool(arr []interface{}) []bool {
	res := make([]bool, len(arr))
	for i, value := range arr {
		b := value.(bool)
		res[i] = b
	}
	return res
}

// ConvertFloat64 uses type assertion to convert an interface slice to a float64 slice
func ConvertFloat64(arr []interface{}) []float64 {
	res := make([]float64, len(arr))
	for i, value := range arr {
		b := value.(float64)
		res[i] = b
	}
	return res
}

// ConvertInt uses type assertion to convert an interface slice to an int slice
func ConvertInt(arr []interface{}) []int {
	res := make([]int, len(arr))
	for i, value := range arr {
		b := value.(int)
		res[i] = b
	}
	return res
}
