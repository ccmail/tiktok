package util

import (
	"fmt"
	"testing"
)

func TestSpliceKey(t *testing.T) {
	key := SpliceKey("aaa", 111, 222, 333)
	fmt.Println(key)
}
