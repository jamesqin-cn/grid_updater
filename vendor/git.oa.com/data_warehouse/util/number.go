package util

import (
	"strconv"
	"strings"
)

var thousand_sep = byte(',')
var dec_sep = byte('.')

func AnyToFloat(val interface{}) float64 {
	var v float64
	switch val.(type) {
	case int:
		v = float64(val.(int))
	case float32:
		v = float64(val.(float32))
	case float64:
		v = val.(float64)
	case string:
		v, _ = strconv.ParseFloat(val.(string), 64)
	default:
		v = float64(0)
	}
	return v
}

// Parses the given float64 into a string, using thousand_sep to separate
// each block of thousands before the decimal point character.
func NumberFormat(val interface{}, precision int) string {
	// Parse the float as a string, with no exponent, and keeping precision
	// number of decimal places. Note that the precision passed in to FormatFloat
	// must be a positive number.
	use_precision := precision
	if precision < 1 {
		use_precision = 1
	}

	as_string := strconv.FormatFloat(AnyToFloat(val), 'f', use_precision, 64)

	// Split the string at the decimal point separator.
	separated := strings.Split(as_string, ".")

	before_decimal := separated[0]
	// Our final string will need a total space of the original parsed string
	// plus space for an additional separator character every 3rd character
	// before the decimal point.
	with_separator := make([]byte, 0, len(as_string)+(len(before_decimal)/3))

	// Deal with a (possible) negative sign:
	if before_decimal[0] == '-' {
		with_separator = append(with_separator, '-')
		before_decimal = before_decimal[1:]
	}

	// Drain the initial characters that are "left over" after dividing the length
	// by 3. For example, if we had "12345", this would drain "12" from the string
	// append the separator character, and ensure we're left with something
	// that is exactly divisible by 3.
	initial := len(before_decimal) % 3
	if initial > 0 {
		with_separator = append(with_separator, before_decimal[0:initial]...)
		before_decimal = before_decimal[initial:]
		if len(before_decimal) >= 3 {
			with_separator = append(with_separator, thousand_sep)
		}
	}

	// For each chunk of 3, append it and add a thousands separator,
	// slicing off the chunks of 3 as we go.
	for len(before_decimal) >= 3 {
		with_separator = append(with_separator, before_decimal[0:3]...)
		before_decimal = before_decimal[3:]
		if len(before_decimal) >= 3 {
			with_separator = append(with_separator, thousand_sep)
		}
	}
	// Append everything after the '.', but only if we have positive precision.
	if precision > 0 {
		with_separator = append(with_separator, dec_sep)
		with_separator = append(with_separator, separated[1]...)
	}
	return string(with_separator)
}
