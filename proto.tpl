syntax = "proto3";

package {{.Package}}

option go_package = "{{.GoPackage}}";

message {{.Name}} {
{{range .MessageDetail}} {{.Type}} {{.Name}} = {{.Num}} [json_name = "{{.Name}}"];
{{end}}}
