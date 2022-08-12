package main


type Message struct {
	Name          string
	Package string
	MessageDetail []TableFied
}
type Field struct {
	TypeName string
	AttrName string
	Num      int
}

type TableFied struct {
	Name string
	Type string
	Comment string
	Num int
}

var typeArr = map[string]string{
	"int":       "int32",
	"tinyint":   "int32",
	"smallint":  "int32",
	"mediumint": "int32",
	"enum":      "int32",
	"bigint":    "int64",
	"varchar":   "string",
	"timestamp": "string",
	"date":      "string",
	"text":      "string",
	"double":    "double",
	"decimal":   "double",
	"float":     "float",
}