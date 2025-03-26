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
	"fmt"

	"google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  fgainterceptors "go.linka.cloud/go-openfga/interceptors"
)

var (
	_ = codes.OK
	_ = status.New
	_ = fmt.Sprintf
	_ = context.Canceled
)

{{ range .Services }}

{{ $types := (types .) }}
const (
{{- range $types }}
	{{ $name := .Name }}
	FGA{{ upperCamelCase .Name }}Type = "{{ $name }}"
	{{ range $key, $value := .Relations }}
	FGA{{ upperCamelCase (printf "%s_%s" $name $key) }} = "{{ $key }}"
	{{- end }}
{{- end }}
)
{{- range $types }}
	{{ $name := .Name }}
// FGA{{ upperCamelCase $name }}Object returns the object string for the {{ $name }} type, e.g. "{{ $name }}:id"
func FGA{{ upperCamelCase $name }}Object(id string) string {
	return FGA{{ upperCamelCase $name }}Type + ":" + id
}
{{- end }}

{{ $service := . }}
{{ with (module .) }}
//go:embed {{ file $service ".fga" }}
var FGAModel string
{{ end }}

// RegisterFGA registers the {{ $service.Name }} service with the provided FGA interceptors.
func RegisterFGA(fga fgainterceptors.FGA) {
	{{- range .Methods }}
	  {{- $method := . }}
	  {{- with access . }}
		{{- $fullMethodName := (printf "%s_%s_FullMethodName" $service.Name $method.Name) }}
	  fga.Register({{ $fullMethodName }}, func(ctx context.Context, req any, user string, kvs ...any) error {
			{{- range .Check }}
			{
				{{- if (need_getter .) }}
				r, ok := req.(*{{ name $method.Input }})
				if !ok {
					panic("unexpected request type: expected {{ name $method.Input }}")
				}
				id := r.{{ getter . $method }}()
				if id == "" {
					return status.Error(codes.InvalidArgument, "{{ field . }} is required")
				}
				object := "{{ .GetType }}" + ":" + fga.Normalize(id)
				{{- if not .GetIgnoreNotFound }}
				ok, err := fga.Has(ctx, object)
				if err != nil {
					return status.Errorf(codes.Internal, "permission check failed: %v", err)
				}
				if !ok {
					return status.Errorf(codes.NotFound, "{{ .GetType }} %q not found", id)
				}
				{{- end }}
				msg := fmt.Sprintf("[%s]: not allowed to call %s on {{ .GetType }} %q", user, {{ $fullMethodName }}, id)
				{{- else }}
				object := "{{ .GetType }}" + ":" + fga.Normalize("{{ .GetID }}")
				msg := fmt.Sprintf("[%s]: not allowed to call %s", user, {{ $fullMethodName }})
				{{- end }}
				granted, err := fga.Check(ctx, object, FGA{{ upperCamelCase (printf "%s_%s" .GetType .GetCheck) }}, user, kvs...)
				if err != nil {
					return status.Errorf(codes.Internal, "permission check failed: %v", err)
				}
				if !granted {
					return status.Error(codes.PermissionDenied, msg)
				}
				return nil
			}
			{{- end }}
    })
	  {{- end }}
	{{- end }}
}
{{ end }}
