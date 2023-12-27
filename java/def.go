package java

import "C"
import (
	"fmt"
	"github.com/MegaXChan/gojni/utils"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/MegaXChan/gojni/jni"
)

var _is64Bit *bool = nil

// 判断系统是不是64位系统
func is64Bit() bool {
	if _is64Bit == nil {
		arch := runtime.GOARCH
		b := strings.Contains(arch, "64")
		_is64Bit = new(bool)
		*_is64Bit = b
	}
	return *_is64Bit
}

// 处理真实c语言类型的类型长度
func sizeOf(kind reflect.Kind, index int) uintptr {
	if jni.ISDEBUG {
		fmt.Println("kind", kind)
	}
	//发现从第9个参数开始才开始使用padding
	if index < 6 {
		if is64Bit() {
			return 8
		} else {
			return 4
		}
	}
	switch kind {
	case reflect.Int, reflect.Uint32, reflect.Float32:
		return 4
	case reflect.Int8, reflect.Uint8, reflect.Bool:
		return 1
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int64, reflect.Uint64:
		return 8
	default:
		//todo 没验证32位系统
		if is64Bit() {
			return 8
		} else {
			return 4
		}
	}
}

func FixPadding(p unsafe.Pointer, thiskind reflect.Kind) unsafe.Pointer {
	arch := runtime.GOARCH
	fix := uintptr(4)
	if strings.Contains(arch, "64") {
		fix = uintptr(8)
	}
	px := uintptr(p)
	mod := px % fix
	if mod != 0 {
		if jni.ISDEBUG {
			fmt.Println("old point", px)
		}
		px2 := px + fix - mod
		thiskindsize := sizeOf(thiskind, 100)
		if px+thiskindsize < px2 {
			if jni.ISDEBUG {
				fmt.Println("new point", px)
			}
			return unsafe.Pointer(px)
		} else {
			if jni.ISDEBUG {
				fmt.Println("new point", px2)
			}
			return unsafe.Pointer(px2)
		}
	}
	return unsafe.Pointer(px)
}

func router2(s string, p0, p1 uintptr, p *uintptr) uintptr {
	defer func() {
		if r := recover(); r != nil {
			if str, b := r.(string); b {
				jni.ThrowException(str)
				return
			}
			if e, b := r.(error); b {
				jni.ThrowException(e.Error())
				return
			}
		}
	}()
	if f, b := fMappers[s]; b {
		setSelfClassOrObject(p1)
		defer clearSelfClassOrObject()

		rValues := reflect.ValueOf(f.fn).Call(convertParam2(f, p))
		if len(rValues) != 1 {
			return 0
		}
		//convert return values
		return utils.JabValueToUint(rValues[0])
	}
	return 0
}

// in c type
func router(s string, p ...uintptr) uintptr {
	defer func() {
		if r := recover(); r != nil {
			if str, b := r.(string); b {
				jni.ThrowException(str)
				return
			}
			if e, b := r.(error); b {
				jni.ThrowException(e.Error())
				return
			}
		}
	}()
	if f, b := fMappers[s]; b {
		setSelfClassOrObject(p[1])
		defer clearSelfClassOrObject()
		rValues := reflect.ValueOf(f.fn).Call(convertParam(f, p...))
		if len(rValues) != 1 {
			return 0
		}
		//convert return values
		return utils.JabValueToUint(rValues[0])
	}
	return 0
}

func convertParam2(f method, params *uintptr) []reflect.Value {
	pointer := unsafe.Pointer(params)
	var ret []reflect.Value
	lenP := len(f.sig)
	env := jni.AutoGetCurrentThreadEnv()
	for i := 0; i < lenP; i++ {
		s := f.sig[i]
		pointer = FixPadding(pointer, s.gSig.Kind())
		switch s.gSig.Kind() {
		case reflect.Uintptr:
			val := *(*uintptr)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Int8:
			val := *(*int8)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Int16:
			val := *(*int16)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Int:
			val := *(*int)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Int32:
			val := *(*int32)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Float32:
			//todo 这个好像不行
			val := *(*float32)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Float64:
			//todo 这个好像不行
			val := *(*float64)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Int64:
			val := *(*int64)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Uint:
			val := *(*uint)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Uint8:
			val := *(*uint8)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Uint16:
			val := *(*uint16)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Uint32:
			val := *(*uint32)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Uint64:
			val := *(*uint64)(pointer)
			ret = append(ret, reflect.ValueOf(val))
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			if jni.ISDEBUG {
				fmt.Println("pointer ->", i, pointer, val)
			}
		case reflect.Bool:
			val := *(*uint8)(pointer)
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))

			if val == jni.JNI_TRUE {
				ret = append(ret, reflect.ValueOf(true))
			} else if val == jni.JNI_FALSE {
				ret = append(ret, reflect.ValueOf(false))
			} else {
				panic("unknown bool")
			}
			//ret = append(ret, reflect.ValueOf(false))
		case reflect.String:
			val := *(*uintptr)(pointer)
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))

			jni.CheckNull(val, "jni input str is null")
			pkg := string(env.GetStringUTF(val))
			ret = append(ret, reflect.ValueOf(pkg))
		case reflect.Slice:
			val := *(*uintptr)(pointer)
			pointer = unsafe.Add(pointer, sizeOf(reflect.ValueOf(val).Kind(), i))
			jni.CheckNull(val, "jni input slice is null")
			ret = append(ret, convertParamSlice(env, s.gSig, val))
		default:
			panic(fmt.Sprintf("err convertParam %v not support", s.gSig.Kind()))
		}
	}
	return ret
}

func convertParam(f method, params ...uintptr) []reflect.Value {
	var ret []reflect.Value
	lenP := len(f.sig)
	env := jni.AutoGetCurrentThreadEnv()
	for i := 0; i < lenP; i++ {
		s := f.sig[i]
		p := params[i+2]
		if jni.ISDEBUG {
			fmt.Println("val -> ", i, p)
		}
		switch s.gSig.Kind() {
		case reflect.Uintptr:
			ret = append(ret, reflect.ValueOf(uintptr(p)))
		case reflect.Int8:
			ret = append(ret, reflect.ValueOf(int8(p)))
		case reflect.Int16:
			ret = append(ret, reflect.ValueOf(int16(p)))
		case reflect.Int:
			ret = append(ret, reflect.ValueOf(int(p)))
		case reflect.Int32:
			ret = append(ret, reflect.ValueOf(int32(p)))
		case reflect.Float32:
			ret = append(ret, reflect.ValueOf(float32(p)))
		case reflect.Float64:
			ret = append(ret, reflect.ValueOf(float64(p)))
		case reflect.Int64:
			ret = append(ret, reflect.ValueOf(int64(p)))
		case reflect.Uint:
			ret = append(ret, reflect.ValueOf(uint(p)))
		case reflect.Uint8:
			ret = append(ret, reflect.ValueOf(uint8(p)))
		case reflect.Uint16:
			ret = append(ret, reflect.ValueOf(uint16(p)))
		case reflect.Uint32:
			ret = append(ret, reflect.ValueOf(uint32(p)))
		case reflect.Uint64:
			ret = append(ret, reflect.ValueOf(uint64(p)))
		case reflect.Bool:

			if p == jni.JNI_TRUE {
				ret = append(ret, reflect.ValueOf(true))
			} else if p == jni.JNI_FALSE {
				ret = append(ret, reflect.ValueOf(false))
			} else {
				panic("unknown bool")
			}
			//ret = append(ret, reflect.ValueOf(false))
		case reflect.String:
			jni.CheckNull(p, "jni input str is null")
			pkg := string(env.GetStringUTF(p))
			ret = append(ret, reflect.ValueOf(pkg))
		case reflect.Slice:
			jni.CheckNull(p, "jni input slice is null")
			ret = append(ret, convertParamSlice(env, s.gSig, p))
		default:
			panic(fmt.Sprintf("err convertParam %v not support", s.gSig.Kind()))
		}
	}
	return ret
}

func convertParamSlice(env *jni.Env, Array reflect.Type, p uintptr) reflect.Value {

	iLen := env.GetArrayLength(p)
	item := Array.Elem()

	switch item.Kind() {
	case reflect.Int:
		iTypes := int(unsafe.Sizeof(C.long(0)))
		jBytes := iLen * iTypes
		ptr := env.GetLongArrayElements(p, true)
		reBytes := C.GoBytes(ptr, C.int(jBytes))
		env.ReleaseLongArrayElements(p, uintptr(ptr), 0)
		head := (*reflect.SliceHeader)(unsafe.Pointer(&reBytes))
		head.Cap /= iTypes
		head.Len /= iTypes
		return reflect.ValueOf(*(*[]int)(unsafe.Pointer(head)))
	case reflect.Int32:
		iTypes := int(unsafe.Sizeof(C.int(0)))
		jBytes := iLen * iTypes
		ptr := env.GetIntArrayElements(p, true)
		reBytes := C.GoBytes(ptr, C.int(jBytes))
		env.ReleaseIntArrayElements(p, uintptr(ptr), 0)
		head := (*reflect.SliceHeader)(unsafe.Pointer(&reBytes))
		head.Cap /= iTypes
		head.Len /= iTypes
		return reflect.ValueOf(*(*[]int32)(unsafe.Pointer(head)))
	case reflect.String:
		var temp = make([]string, iLen)
		for i := 0; i < iLen; i++ {
			temp[i] = string(env.GetStringUTF(env.GetObjectArrayElement(p, i)))
		}
		return reflect.ValueOf(temp)
	case reflect.Uint8:
		iTypes := 1
		jBytes := iLen * iTypes
		ptr := env.GetByteArrayElements(p, true)
		reBytes := C.GoBytes(ptr, C.int(jBytes))
		env.ReleaseByteArrayElements(p, uintptr(ptr), 0)
		head := (*reflect.SliceHeader)(unsafe.Pointer(&reBytes))
		head.Cap /= iTypes
		head.Len /= iTypes
		return reflect.ValueOf(*(*[]byte)(unsafe.Pointer(head)))
	case reflect.Float32:
		iTypes := int(unsafe.Sizeof(C.float(0.0)))
		jBytes := iLen * iTypes
		ptr := env.GetFloatArrayElements(p, true)
		reBytes := C.GoBytes(ptr, C.int(jBytes))
		env.ReleaseFloatArrayElements(p, uintptr(ptr), 0)
		head := (*reflect.SliceHeader)(unsafe.Pointer(&reBytes))
		head.Cap /= iTypes
		head.Len /= iTypes
		return reflect.ValueOf(*(*[]float32)(unsafe.Pointer(head)))
	case reflect.Float64:
		iTypes := int(unsafe.Sizeof(C.double(0.0)))
		jBytes := iLen * iTypes
		ptr := env.GetDoubleArrayElements(p, true)
		reBytes := C.GoBytes(ptr, C.int(jBytes))
		env.ReleaseDoubleArrayElements(p, uintptr(ptr), 0)
		head := (*reflect.SliceHeader)(unsafe.Pointer(&reBytes))
		head.Cap /= iTypes
		head.Len /= iTypes
		return reflect.ValueOf(*(*[]float64)(unsafe.Pointer(head)))
	default:
		panic(fmt.Sprintf("not support Array %s ", item))
	}
	return reflect.Value{}
}
