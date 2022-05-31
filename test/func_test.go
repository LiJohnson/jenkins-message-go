package test

import (
	"fmt"
	"testing"
)

func TestAAA(t *testing.T) {
	fmt.Println("testing")
}

func Test_array(t *testing.T) {
	var list [123]int
	fmt.Println(len(list), cap(list))
}
