package utils

import (
	"strconv"
	"strings"
)

func StringToInt8Array(s string) ([]int8, error) {
	result := make([]int8, len(s))

	arrStr := strings.Split(s, "")
	for i, v := range arrStr {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		result[i] = int8(n)
	}
	return result, nil
}
