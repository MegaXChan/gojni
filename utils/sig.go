package utils

//jni sig
var SigEncodeMap = map[string]string{
	"int":     "I",
	"boolean": "Z",
	"byte":    "B",
	"char":    "C",
	"short":   "S",
	"long":    "J",
	"float":   "F",
	"double":  "D",
	"void":    "V",
}

var SigDecodeMap = map[string]string{
	"I": "int",
	"Z": "boolean",
	"B": "byte",
	"C": "char",
	"S": "short",
	"J": "long",
	"F": "float",
	"D": "double",
	"V": "void",
}

/*
	stop ()V
	args ([I)V
	nice ([Ljava/lang/String;)V
	bb ([B)V
	llll ([J)V
	fff ([F)V
	ddd ([D)V
*/
func SigToJavaNative(name string, sig string) {

}
