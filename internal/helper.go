package internal

import (
	"fmt"
)

func LogString(v any) string {
	return fmt.Sprintf("%+v", v)
}
