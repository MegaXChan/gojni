package java

import (
	"fmt"
	"runtime"
	"testing"
	"unsafe"
)

func TestCode(t *testing.T) {
	fmt.Println(unsafe.Sizeof(float32(0.0)))
}

func TestNum(f *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if ok {

		fmt.Println(filename)

	}
	for i := 'A'; i < 'A'+26+26; i++ {
		fmt.Printf("%c\n", i)
	}
	//fmt.Sprintf("%c%d", dep+97, inNum)
}
