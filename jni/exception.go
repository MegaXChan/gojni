package jni

import (
	"fmt"
	"runtime/debug"
)

var ISDEBUG = false

type JavaExceptionCodes int

const (
	JavaOutOfMemoryError JavaExceptionCodes = iota
	JavaIOException
	JavaRuntimeException
	JavaIndexOutOfBoundsException
	JavaArithmeticException
	JavaIllegalArgumentException
	JavaNullPointerException
	JavaDirectorPureVirtual
	JavaUnknownError
	JavaException
)

var exceptionMap = map[JavaExceptionCodes]string{
	JavaOutOfMemoryError:          "java/lang/OutOfMemoryError",
	JavaIOException:               "java/io/IOException",
	JavaRuntimeException:          "java/lang/RuntimeException",
	JavaIndexOutOfBoundsException: "java/lang/IndexOutOfBoundsException",
	JavaArithmeticException:       "java/lang/ArithmeticException",
	JavaIllegalArgumentException:  "java/lang/ArithmeticException",
	JavaNullPointerException:      "java/lang/NullPointerException",
	JavaDirectorPureVirtual:       "java/lang/RuntimeException",
	JavaUnknownError:              "java/lang/UnknownError",
	JavaException:                 "java/lang/Exception",
}

func ThrowException(msg string) {
	if ISDEBUG {
		NativeThrowException(AutoGetCurrentThreadEnv(), JavaException, msg+"\n"+string(debug.Stack()))
	} else {
		NativeThrowException(AutoGetCurrentThreadEnv(), JavaException, msg)
	}

}

func JavaThrowException(msg string) {
	env := AutoGetCurrentThreadEnv()
	jCls := env.FindClass("java/lang/Exception")
	if jCls == 0 {
		panic("not find java/lang/Exception")
	}
	jMethodID := env.GetMethodID(jCls, "<init>", "(Ljava/lang/String;)V")
	job := env.NewObjectA(jCls, jMethodID, Jvalue(env.NewString(msg)))
	jPrint := env.GetMethodID(jCls, "printStackTrace", "()V")
	env.CallVoidMethodA(job, jPrint)
}

func NativeThrowException(env *Env, code JavaExceptionCodes, msg string) {
	CheckException(env)
	if cls, b := exceptionMap[code]; b {
		jcls := env.FindClass(cls)
		ok := env.ThrowNew(jcls, msg+" [from native]")
		fmt.Println(ok)
	}
}

func ExceptionMessageFromThrowable(JNIEnv Env, jthrowable Jthrowable) Jstring {
	JNIEnv.ExceptionClear()
	throwclz := JNIEnv.GetObjectClass(jthrowable)
	if throwclz > 0 {
		getMessageMethodID := JNIEnv.GetMethodID(throwclz, "getMessage", "()Ljava/lang/String;")
		if getMessageMethodID > 0 {
			jmsg := JNIEnv.CallObjectMethodA(jthrowable, getMessageMethodID)
			return jmsg
		}
	}
	return 0
}

func PrintException(JNIEnv *Env, jthrowable Jthrowable) {
	if jthrowable == 0 {
		return
	}
	throwclz := JNIEnv.FindClass("java/lang/Throwable")
	printStackMethod := JNIEnv.GetMethodID(throwclz, "printStackTrace", "()V")
	JNIEnv.CallNonvirtualVoidMethodA(jthrowable, throwclz, printStackMethod)
	JNIEnv.ExceptionClear()
}

/**
 *
 */
func CheckNullException(msg string, ok func(env *Env), checkNull ...uintptr) {
	has := false
	s := ""
	for i, u := range checkNull {
		if u == 0 {
			s = fmt.Sprintf(" [ check args %d is null ]", i)
			has = true
			break
		}
	}
	env := AutoGetCurrentThreadEnv()
	if has {
		NativeThrowException(env, JavaNullPointerException, msg+s)
	} else {
		ok(env)
	}
}

func CheckNull(uin uintptr, msg string) {
	if uin == 0 {
		env := AutoGetCurrentThreadEnv()
		NativeThrowException(env, JavaNullPointerException, msg)
	}
}

func CheckException(env *Env) {
	if env.ExceptionCheck() {
		PrintException(env, env.ExceptionOccurred())
		panic("CheckException")
	}
}
