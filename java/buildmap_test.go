package java

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestBuildFile(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fail()
		return
	}
	file := path.Dir(filename) + "/method_def.go"

	maxArgs := 16
	maxMethod := 256

	builder := strings.Builder{}
	builder.WriteString("package java")
	builder.WriteString("\n\n")

	methodstr := "\n//extern void* m_%v_%v(%s);"
	for i := 2; i <= maxArgs; i++ {
		for j := 0; j < maxMethod; j++ {
			list := []string{}
			for k := 0; k <= i; k++ {
				list = append(list, "void *")
			}
			builder.WriteString(fmt.Sprintf(methodstr, j, i, strings.Join(list, ",")))
		}
	}
	builder.WriteString("\nimport \"C\"")
	builder.WriteString(fmt.Sprintf(`
//import (
//	"unsafe"
//)
`))

	builder.WriteString("\nvar nMap = FuncMap{")
	for i := 2; i <= maxArgs; i++ {
		list := []string{}
		for j := 0; j < maxMethod; j++ {
			methodName := fmt.Sprintf("funcc{code: \"m_%v_%v\", fun: C.m_%v_%v}", j, i, j, i)
			list = append(list, methodName)
		}
		builder.WriteString(fmt.Sprintf("\n\t\t%v:  {%s},", i, strings.Join(list, ",")))
	}

	builder.WriteString("\n}")

	os.WriteFile(file, []byte(builder.String()), 0666)
	//-------------------------

	file2 := path.Dir(filename) + "/method_init.go"

	builder = strings.Builder{}
	builder.WriteString("package java")
	builder.WriteString("\n\n")
	builder.WriteString("\nimport \"C\"")

	builder.WriteString("\n")

	builder.WriteString(fmt.Sprintf(`
const (
	maxArgs = %v
	maxMethod   = %v
)`, maxArgs, maxMethod))

	for i := 2; i <= maxArgs; i++ {
		for j := 0; j < maxMethod; j++ {
			list := []string{}
			for k := 0; k <= i; k++ {
				list = append(list, fmt.Sprintf("p%v", k))
			}
			methodName := fmt.Sprintf("m_%v_%v", j, i)
			argss := strings.Join(list, ",")
			builder.WriteString(fmt.Sprintf("\n\n//export %v", methodName))

			builder.WriteString(fmt.Sprintf(`
func %v(%v uintptr) uintptr {
	return router("%v", %v)
}`, methodName, argss, methodName, argss))

		}
	}

	os.WriteFile(file2, []byte(builder.String()), 0666)

	//for i := 'A'; i < 'A'+26+26; i++ {
	//	fmt.Printf("%c\n", i)
	//}
	//fmt.Sprintf("%c%d", dep+97, inNum)
}

func TestBuildFileExt(t *testing.T) {
	//for i := 0; i < 100; i++ {
	//	fmt.Println(genMethodName())
	//}

	maxArgs := 16
	maxMethod := 256

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fail()
		return
	}

	type funT struct {
		name   string
		arglen int
	}

	nmap := map[int][]funT{}

	for i := 2; i <= maxArgs; i++ {
		list := []funT{}
		for j := 0; j < maxMethod; j++ {
			f := funT{
				name:   genMethodName(),
				arglen: i,
			}
			list = append(list, f)
		}
		nmap[i] = list
	}

	filedef := path.Dir(filename) + "/method_def.go"
	fileinit := path.Dir(filename) + "/method_init.go"

	builder := strings.Builder{}

	builder.WriteString("package java")
	builder.WriteString("\n\n")

	methodstr := "\n//extern void* %s(%s);"
	for _, v := range nmap {
		for _, f := range v {
			name := f.name
			arglen := f.arglen
			list := []string{}
			for j := 0; j < arglen; j++ {
				list = append(list, "void *")
			}

			builder.WriteString(fmt.Sprintf(methodstr, name, strings.Join(list, ",")))
		}
	}

	builder.WriteString("\nimport \"C\"")
	builder.WriteString(fmt.Sprintf(`
//import (
//	"unsafe"
//)
`))

	builder.WriteString("\nvar nMap = FuncMap{")

	for k, v := range nmap {
		list := []string{}
		for _, f := range v {
			name := f.name
			methodName := fmt.Sprintf("funcc{code: \"%s\", fun: C.%s}", name, name)
			list = append(list, methodName)
		}
		builder.WriteString(fmt.Sprintf("\n\t\t%v:  {%s},", k, strings.Join(list, ",")))
	}

	builder.WriteString("\n}")

	os.WriteFile(filedef, []byte(builder.String()), 0666)
	//-------------------------

	builder = strings.Builder{}
	builder.WriteString("package java")
	builder.WriteString("\n\n")
	builder.WriteString("\nimport \"C\"")

	builder.WriteString("\n")

	builder.WriteString(fmt.Sprintf(`
const (
	maxArgs = %v
	maxMethod   = %v
)`, maxArgs, maxMethod))

	for _, v := range nmap {

		for _, f := range v {
			name := f.name
			arglen := f.arglen
			list := []string{}
			for k := 0; k <= arglen; k++ {
				list = append(list, fmt.Sprintf("p%v", k))
			}

			argss := strings.Join(list, ",")
			builder.WriteString(fmt.Sprintf("\n\n//export %v", name))

			builder.WriteString(fmt.Sprintf(`
func %v(%v uintptr) uintptr {
	return router("%v", %v)
}`, name, argss, name, argss))

		}

	}

	os.WriteFile(fileinit, []byte(builder.String()), 0666)

}

var methodNameMap = map[string]int{}

func genMethodName() string {
	methodlen := 8
	methodNameStart := []string{}
	methodNameVal := []string{}

	for i := 'a'; i <= 'z'; i++ {
		methodNameStart = append(methodNameStart, fmt.Sprintf("%c", i))
		methodNameVal = append(methodNameVal, fmt.Sprintf("%c", i))
	}
	for i := 'A'; i <= 'Z'; i++ {
		methodNameStart = append(methodNameStart, fmt.Sprintf("%c", i))
		methodNameVal = append(methodNameVal, fmt.Sprintf("%c", i))
	}
	for i := '0'; i <= '9'; i++ {
		methodNameVal = append(methodNameVal, fmt.Sprintf("%c", i))
	}
	method := []string{}
	for i := 0; i < methodlen; i++ {
		rad := int(math.Abs(float64(rand.Int())))
		if i == 0 {
			index := rad % len(methodNameStart)
			method = append(method, methodNameStart[index])
		} else {
			index := rad % len(methodNameVal)
			method = append(method, methodNameVal[index])
		}
	}
	name := strings.Join(method, "")
	fmt.Println("method Name ->", name)
	_, ok := methodNameMap[name]
	if ok {
		fmt.Println("got same name")
		return genMethodName()
	} else {
		methodNameMap[name] = 0
	}
	return name
}
