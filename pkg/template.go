package lp

// https://primer.style/css/getting-started
// https://primer.style/css/utilities/flexbox
// https://primer.style/css/components/select-menu

const commonTemplate = `<!doctype html>
<html data-color-mode="dark" data-dark-theme="dark_dimmed" lang="en">
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
        <a class="UnderlineNav-item" href="{{ .Name }}.html">{{ .Name }}</a>
      {{ end }}
      <a class="UnderlineNav-item" href="https://github.com/rjbrown57/lp.git">lp</a>
    </div>
    </div>
  </div>
</nav>
`

const bodyTemplate = `
<body>
    <div class="m-2 d-flex color-bg-subtle flex-justify-around">
      {{- range .Headings}}
        <div class="p-5 border rounded color-border-danger color-shadow-medium">
        <h3>{{ .Name }}</h3>
        <p>
          <ul>
          {{- range .Links }}
          <li><a href="{{ .Url }}">{{ .Name }}</a></li>
          {{- end }}
          <ul>
        </p>
        </div>
      {{- end }}
    </div>
</body>
</html>
`
