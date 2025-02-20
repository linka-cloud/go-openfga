{{- /*gotype: github.com/lyft/protoc-gen-star.File*/ -}}
{{ comment .SyntaxSourceCodeInfo.LeadingComments }}
{{ range .SyntaxSourceCodeInfo.LeadingDetachedComments }}
{{ comment . }}
{{ end }}
// Code generated by protoc-gen-go-openfga. DO NOT EDIT.
package {{ package . }}

import (
  "context"
	_ "embed"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  fgainterceptors "go.linka.cloud/go-openfga/interceptors"
)

var (
	_ = codes.OK
	_ = status.New
)

{{ range .Services }}
{{ $service := . }}
{{ with (module .) }}
//go:embed {{ file $service ".fga" }}
var FGAModel string

const (
	{{- range .Extends }}
	{{- $name := .Type }}
	{{- with .Relations }}
	FGA{{ upperCamelCase $name }}Type = "{{ $name }}"
	{{- range . }}
	FGA{{ upperCamelCase (printf "%s_%s" $name .Define) }} = "{{ .Define }}"
	{{- end }}
	{{- end }}
	{{- end }}
	{{- range .Definitions }}
	{{- $name := .Type }}
	{{- with .Relations }}
	FGA{{ upperCamelCase $name }}Type = "{{ $name }}"
	{{- range . }}
	FGA{{ upperCamelCase (printf "%s_%s" $name .Define) }} = "{{ .Define }}"
	{{- end }}
	{{- end }}
	{{- end }}
)

{{- range .Extends }}
{{- $name := .Type }}
{{- with .Relations }}
func FGA{{ upperCamelCase $name }}Object(id string) string {
	return FGA{{ upperCamelCase $name }}Type + ":" + id
}
{{- end }}
{{- end }}
{{- range .Definitions }}
{{- $name := .Type }}
{{- with .Relations }}
func FGA{{ upperCamelCase $name }}Object(id string) string {
	return FGA{{ upperCamelCase $name }}Type + ":" + id
}
{{- end }}
{{- end }}
{{ end }}

func RegisterFGA(fga fgainterceptors.FGA) {
	{{- range .Methods }}
	  {{- $method := . }}
	  {{- with access . }}
	  fga.Register({{ $service.Name }}_{{ $method.Name }}_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
			{{- if (need_getter .) }}
			r, ok := req.(*{{ name $method.Input }})
			if !ok {
				panic("unexpected request type: expected {{ name $method.Input }}")
			}
			id := r.{{ getter . $method }}()
			if id == "" {
				return "", "", status.Error(codes.InvalidArgument, "{{ field . }} is required")
			}
			return FGA{{ upperCamelCase .Type }}Object(id), FGA{{ upperCamelCase (printf "%s_%s" .Type .Check) }}, nil
			{{- else }}
			return FGA{{ upperCamelCase .Type }}Type + ":" + "{{ .ID }}", FGA{{ upperCamelCase (printf "%s_%s" .Type .Check) }}, nil
			{{- end }}
    })
	  {{- end }}
	{{- end }}
}
{{ end }}
