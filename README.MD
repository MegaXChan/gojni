## quick start
```shell
go get github.com/MegaXChan/gojni
```
### 1.java invoke native
in go 
```go
package main
import (
	"fmt"
	"unsafe"
	"github.com/MegaXChan/gojni/java"
	"github.com/MegaXChan/gojni/jni"
)
func main() {}
func init() {
	java.OnMainLoad(func(reg java.Register) {
		reg.WithClass("com.nk.Hello").
			BindNative("nice", "void(java.lang.String[])", nice).
			BindNative("setVal", "void(int)", setVal).
			BindNative("New", "com.nk.Hello(com.nk.Hello)", New).
			Done()
	})
}

func New(main jni.Jobject) jni.Jobject {
	env := jni.AutoGetCurrentThreadEnv()
	class := env.GetObjectClass(main)
	aid := env.GetFieldID(class, "a", "I")
	ival := env.GetIntField(main, aid)
	initMethod := env.GetMethodID(class, "<init>", "()V")
	newObj := env.NewObjectA(class, initMethod)

	env.SetIntField(newObj, aid, ival)
	return newObj
}

func nice(ss []string) {
	selfclass := java.GetSelfClassOrObject()
	env := jni.AutoGetCurrentThreadEnv()
	class := env.FindClass("java/lang/Class")
	toString := env.GetMethodID(class, "toString", "()Ljava/lang/String;")
	toStringVal := env.CallObjectMethodA(selfclass, toString)
	pkg := string(env.GetStringUTF(toStringVal))
	fmt.Println(pkg)
	fmt.Println(ss)
}

func setVal(a int32) {
	selfobj := java.GetSelfClassOrObject()
	env := jni.AutoGetCurrentThreadEnv()
	class := env.GetObjectClass(selfobj)
	aid := env.GetFieldID(class, "a", "I")
	env.SetIntField(selfobj, aid, int(a))
}
```
in java
```java
package com.nk;

public class Hello {
    static {
        System.loadLibrary("test");
    }
    public static void main(String[] args) {
        nice(new String[] { "come", "from", "gojni" });
        Hello h = new Hello();
        h.setVal(2)
        System.out.println(h.a);
        
        Hello nh = h.New(h);
        System.out.println(nh.a);
    }
    public static native void nice(String[] sss);
    public native void setVal(int val);
    native Hello New(Hello x);
    public int a;
}
```
### build and run
```
go build -buildmode=c-shared -o libtest.so
LD_LIBRARY_PATH=. java com.nk.Hello
```

### 2.native invoke java
```go
native.LoadClass("com.szb.Jvm").New().Invoke("attach", "void()")

fmt.Println(vm.RunSource("java.lang.Thread.currentThread[java.lang.Thread()]();").
				AsObject().
				Invoke("getName", "java.lang.String()").
				AsString())
```
# Android build
```
GOOS=android \
GOARCH=arm64 \
CC=$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang \
CXX=$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang++ \
CGO_ENABLED=1 \
go build -buildmode=c-shared -o aaa.so
```

# IOS build
```
GOOS=darwin \
GOARCH=arm64 \
CC=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang \
CXX=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang++ \
CGO_CFLAGS="-isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS13.7.sdk -miphoneos-version-min=7.0 -fembed-bitcode -arch arm64" \
CGO_CXXFLAGS="-isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS13.7.sdk -miphoneos-version-min=7.0 -fembed-bitcode -arch arm64" \
CGO_LDFLAGS="-isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS13.7.sdk -miphoneos-version-min=7.0 -fembed-bitcode -arch arm64" \
CGO_ENABLED=1 \
go build -tags ios -v -x -buildmode=c-archive -o aaa.a
```
