package util

import (
	"fmt"
)

// SpliceKey 用于拼接Redis的Key
// eg: a= Fans. b = 11, c=22 return:"Fans:[11 22]"u
func SpliceKey(a string, b ...any) string {
	return fmt.Sprint(a, ":", b)
}
