package java

import "C"
import (
	"fmt"
	"github.com/MegaXChan/gojni/utils"
	"reflect"
	"unsafe"

	"github.com/MegaXChan/gojni/jni"
)

//type ArgVal struct {
//	Obj uintptr
//}

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
