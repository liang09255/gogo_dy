message {{.MessageName}} {
    {{range $index, $value := .Result}}
        {{$value.Type}} {{$value.Field}} = {{Add $index 1}};
    {{end}}
}