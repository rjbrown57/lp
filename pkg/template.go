package lp

// https://primer.style/css/getting-started
// https://primer.style/css/utilities/flexbox
// https://primer.style/css/components/select-menu

const commonTemplate = `<!doctype html>
<html data-color-mode="dark" data-dark-theme="{{ .Template.Theme }}" lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link href="https://unpkg.com/@primer/css@^16.0.0/dist/primer.css" rel="stylesheet" />
</head>
`

const navbarTemplate = `
<nav class="UnderlineNav UnderlineNav--full">
    <div class="UnderlineNav-body">
      {{- range .Template.Pages }}
        {{- if .IsIndex }}
        <a class="UnderlineNav-item" href="index.html">{{ .Name }}</a>
        {{- else }}
        <a class="UnderlineNav-item" href="{{ .Name }}.html">{{ .Name }}</a>
        {{- end}}
      {{- end }}
      <a class="UnderlineNav-item" href="https://github.com/rjbrown57/lp.git">lp</a>
    </div>
</nav>
`

const bodyTemplate = `
<body>
    <div class="m-2 d-flex color-bg-subtle flex-justify-start flex-wrap">
      {{- range .Headings}}
        <div class="m-4 p-5 border rounded color-border-danger color-shadow-large">
        <b>{{ .Name }}</b>
        <p><ul>
          {{- range .Links }}
          {{- if .Url }}
          <li><a href="{{ .Url }}">{{ .Name }}</a></li>
          {{- end }}
          {{- if .Urls }}
          <div style="m-4">
          <li>{{ .Name }}
          <details class="dropdown details-reset details-overlay d-inline-block">
          <summary class="color-fg-muted p-2 d-inline" aria-haspopup="true">
            <div class="dropdown-caret"></div>
          </summary>
          <ul class="dropdown-menu dropdown-menu-se">
          {{- range $index,$url := .Urls }}
            {{- range $k,$v := $url }}
            <li><a class="dropdown-item" href="{{ $v }}">{{ $k }}</a></li>
            {{- end }}
          {{- end }}
          </ul>
          </details>
          </li>
          </div>
          {{- end }}
          {{- end }}
        </ul></p>
        </div>
      {{- end }}
    </div>
</body>
</html>
`
