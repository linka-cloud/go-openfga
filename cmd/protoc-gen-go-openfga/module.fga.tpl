{{- /*gotype: go.linka.cloud/go-openfga/openfga/v1alpha1.Module */ -}}
# Code generated by protoc-gen-go-openfga. DO NOT EDIT.

module {{ .Name }}
{{ with .Extends }}
{{- range . }}
extend type {{ .Name }}
  {{- with .Relations }}
  relations
    {{- range . }}
    define {{ .Name }}: {{ .Relation }}
    {{- end }}
  {{- end }}
{{- end }}
{{- end }}
{{ with .Types }}

{{- range . }}
type {{ .Name }}
  {{- with .Relations }}
  relations
    {{- range . }}
    define {{ .Name }}: {{ .Relation }}
    {{- end }}
  {{- end }}
{{- end }}
{{- end }}
{{ with .Conditions }}
{{- range . }}
condition {{ . }}
{{- end }}
{{- end }}
