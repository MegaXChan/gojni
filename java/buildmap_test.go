package java

import (
	"fmt"
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
