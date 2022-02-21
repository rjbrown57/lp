<html>
{{- range .Template.Pages }}
<h1>{{ .Name }}</h1>
{{- range .Headings}}
<h3{{ .Name }}</h3>
<p>
{{- range .Links }}
  <a href="{{ .Url }}">{{ .Name }}</a>
{{- end }}
</p>
{{- end }}
{{- end }}
</html>