type {{.StructName}} struct {
    {{range $index,$value := .Result}}
        {{$value.Field}} {{$value.Type}}
    {{end}}
}