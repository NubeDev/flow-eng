package conversions

import (
	"fmt"
)

func FloatToString(f float64) string {
	f = FloatToFixed(f, 2)
	return fmt.Sprintf("%.2f", f)
}
