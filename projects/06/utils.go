package main

import (
	"strconv"
	"strings"
)

// Decimal2Binary transfer a decimal number to binary formate
func Decimal2Binary(n int)string{
	res := ""

	if n == 0 {
		return "000000000000000"
	}

	for ;n > 0;n /= 2 {
		lsb := n % 2
		res = strconv.Itoa(lsb) + res
	}

	// fill the leading zero
	if len(res) <15{
		res= strings.Repeat("0", 15-len(res)) + res
	}
	// truncate low 15-bit
	if len(res) > 15{
		res=res[len(res)-15:]
	}

	return res
}