syntax = "proto3";
option go_package = "{{.Package}}";

message {{.Name}} {
{{range .MessageDetail}} {{.Type}} {{.Name}} = {{.Num}} [json_name = "{{.Name}}"];
{{end}}}
