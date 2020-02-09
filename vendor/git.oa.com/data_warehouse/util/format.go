package util

import (
	"fmt"
)

func RenderColor(text string, val interface{}, lowThreshold float64, hiThreshold float64) string {
	v := AnyToFloat(val)

	if v < lowThreshold {
		return fmt.Sprintf("<span style=\"color:red;\">%s</span>", text)
	}

	if v > hiThreshold {
		return fmt.Sprintf("<span style=\"color:darkgreen;\">%s</span>", text)
	}

	return text
}
