package utils

import (
	"fmt"
	"github.com/MegaXChan/gojni/jni"
	"github.com/petermattis/goid"
	"reflect"
	"sync"
	"testing"
)

func TestSig(t *testing.T) {

}

func _jabValueToUint(r reflect.Value) uintptr {
	switch r.Type().Kind() {
	case reflect.Uintptr:
		return uintptr(r.Uint())
	//case reflect.String:
	//return env.NewString(r.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uintptr(r.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintptr(r.Uint())
	case reflect.Bool:
		if r.Bool() {
			return jni.JNI_TRUE
		}
		return jni.JNI_FALSE
	default:
		panic(fmt.Sprintf("Return not support type %s", r.Kind().String()))
	}
}

func TestJabValueToUint(t *testing.T) {
	var aa uintptr = 100212324

	r := _jabValueToUint(reflect.ValueOf(aa))

	fmt.Println(r)

}

func TestJabValueToUint2(t *testing.T) {
	get := goid.Get()
	fmt.Println(get)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		x := goid.Get()
		fmt.Println("run1:", x)

	}()

	go func() {
		defer wg.Done()
		x := goid.Get()
		fmt.Println("run2:", x)

	}()
	go func() {
		defer wg.Done()
		x := goid.Get()
		fmt.Println("run3:", x)

	}()

	wg.Wait()
}
