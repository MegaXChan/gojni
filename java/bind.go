package java

import "C"
import (
	"fmt"
	"github.com/MegaXChan/gojni/jni"
	"github.com/MegaXChan/gojni/native"
	"github.com/MegaXChan/gojni/utils"
	"reflect"
	"strings"
	"unsafe"
)

type nativeWarp struct {
	sCls         string
	jCls         jni.Jclass
	env          *jni.Env
	natives      []jni.JNINativeMethod
	isRegistered bool
}

type method struct {
	fn  interface{}
	sig []args
}

type args struct {
	jSig string
	gSig reflect.Type
}

type funcc struct {
	code string
	fun  unsafe.Pointer
}

type FuncMap map[int][]funcc

var (
	fMappers = make(map[string]method)
)

type Register interface {
	GetVm() jni.VM
	WithClass(cls string) *nativeWarp
	Done()
}

type registerImpl struct {
	vm       jni.VM
	instance *nativeWarp
}

func SetFuncMap(m FuncMap) {
	nMap = m
}

func (reg *registerImpl) GetVm() jni.VM {
	return reg.vm
}

func (reg *registerImpl) Done() {
	if reg.instance != nil {
		reg.instance.Done()
	}
}

func (reg *registerImpl) WithClass(cls string) *nativeWarp {
	if reg.vm == 0 {
		panic("forbid")
	}
	reg.Done()
	reg.instance = withClass(cls)
	return reg.instance
}

func withClass(cls string) *nativeWarp {
	env := jni.AutoGetCurrentThreadEnv()
	jCls := env.FindClass(strings.ReplaceAll(cls, ".", "/"))
	if jCls == 0 {
		jni.ThrowException(fmt.Sprintf("not find class %s", cls))
	}
	return &nativeWarp{jCls: jCls, sCls: cls, env: env}
}

func (n *nativeWarp) WithClass(cls string) *nativeWarp {
	n.Done()
	return withClass(cls)
}

func (n *nativeWarp) getPFunc(inNum int) funcc {

	funccs, ok := nMap[inNum]
	if !ok {
		panic(fmt.Sprintf("function args overflow numIN %d", inNum))
	}
	i := len(funccs)
	//fmt.Printf("funccs len -> %v,%v\n", inNum, i)
	if i == 0 {
		panic(fmt.Sprintf("function pools overflow"))
	}
	fuc := funccs[0]
	nMap[inNum] = funccs[1:]
	//fmt.Printf("funccs len -> %v,%v\n", inNum, len(nMap[inNum]))

	return fuc
}

func (n *nativeWarp) BindJNINative(javaMethodName string, def string, fun unsafe.Pointer) *nativeWarp {
	jni.CheckNull(n.jCls, fmt.Sprintf("not find class %s", n.sCls))
	ms := native.EncodeToSig(def)
	n.natives = append(n.natives, jni.JNINativeMethod{Name: javaMethodName, Sig: ms.Sig, FnPtr: fun})
	return n
}

func (n *nativeWarp) BindNative(javaMethodName string, def string, fun interface{}) *nativeWarp {
	jni.CheckNull(n.jCls, fmt.Sprintf("not find class %s", n.sCls))
	ms := native.EncodeToSig(def)
	//fmt.Println(ms.Sig)
	inNum := len(ms.ParamTyp) + 2
	goF := reflect.TypeOf(fun)
	if len(ms.ParamTyp) != goF.NumIn() {
		panic(fmt.Sprintf("method %s not match fun %s %d", javaMethodName, ms.ParamTyp, goF.NumIn()))
	}

	var mArgs []args
	for i := 0; i < goF.NumIn(); i++ {
		//n.CheckType(i, javaMethodName, def, ms.ParamTyp[i].GetSigType(), goF.In(i))
		mArgs = append(mArgs, args{
			jSig: ms.ParamTyp[i].GetSigType(),
			gSig: goF.In(i),
		})
	}
	//if goF.NumOut() > 0 {
	//	n.CheckReturn(javaMethodName, ms.RetTyp.GetSigType(), goF.Out(0))
	//}
	f := n.getPFunc(inNum)
	fMappers[f.code] = method{
		fn:  fun,
		sig: mArgs,
	}
	n.natives = append(n.natives, jni.JNINativeMethod{Name: javaMethodName, Sig: ms.Sig, FnPtr: f.fun})
	return n
}

var checkMap = map[string]reflect.Type{
	"[I":                  reflect.TypeOf((*[]int32)(nil)).Elem(),
	"[Ljava/lang/String;": reflect.TypeOf((*[]string)(nil)).Elem(),
	"[B":                  reflect.TypeOf((*[]byte)(nil)).Elem(),
	"[J":                  reflect.TypeOf((*[]int)(nil)).Elem(),
	"[F":                  reflect.TypeOf((*[]float32)(nil)).Elem(),
	"[D":                  reflect.TypeOf((*[]float64)(nil)).Elem(),

	"I":                  reflect.TypeOf((*int32)(nil)).Elem(),
	"Ljava/lang/String;": reflect.TypeOf((*string)(nil)).Elem(),
	"B":                  reflect.TypeOf((*byte)(nil)).Elem(),
	"J":                  reflect.TypeOf((*int)(nil)).Elem(),
	"Z":                  reflect.TypeOf((*bool)(nil)).Elem(),
	"F":                  reflect.TypeOf((*float32)(nil)).Elem(),
	"D":                  reflect.TypeOf((*float64)(nil)).Elem(),
}

func (n *nativeWarp) CheckReturn(mName string, jsig string, gTyp reflect.Type) {
	if v, b := checkMap[jsig]; !b || v != gTyp {
		if b {
			panic(fmt.Sprintf("\n%s method %s return { %s  } not match go type {%s} \nmust use go type ==> %s",
				n.sCls, mName, jsig, gTyp, v))
		} else {
			panic(fmt.Sprintf("%s method %s return { %s  }  not support", n.sCls, mName, jsig))
		}
	}
}

func (n *nativeWarp) CheckType(i int, mName string, def string, jsig string, gTyp reflect.Type) {
	if v, b := checkMap[jsig]; !b || v != gTyp {
		if b {
			panic(fmt.Sprintf("\n%s method %s definition { %s %d } not match go type {%s} \nmust use go type ==> %s",
				n.sCls, mName, def, i, gTyp, v))
		} else {
			panic(fmt.Sprintf("%s method %s definition { %s %d } sig %s not support", n.sCls, mName, def, i, jsig))
		}
	}
}

func (n *nativeWarp) Done() {
	if n.env.RegisterNatives(n.jCls, n.natives) < 0 && !n.isRegistered {
		if jni.ISDEBUG {
			fmt.Println("java class: ", n.sCls)
		}
		n.printNative()
		panic("RegisterNatives error \nplease check java nativeWarp define ")
	} else {
		n.isRegistered = true
	}
}

func (n *nativeWarp) printNative() {
	for _, nativeMethod := range n.natives {
		fmt.Printf("%s %s\n", utils.Wp(nativeMethod.Name, 30), utils.Wp(nativeMethod.Sig, 100))
	}
}
