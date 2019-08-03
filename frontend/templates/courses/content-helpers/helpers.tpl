{{ define "problem" }}
  <h3 class="display">{{ .LocalizedTitle ctx.Locales }}</h3>
  {{ .LocalizedStatement ctx.Locales }}

  {{ range .Samples }}
    <b>Sample</b>
    <pre>{{ .InputFile.FileString }}</pre>
    <hr>
    <pre>{{ .OutputFile.FileString }}</pre>
  {{end}}
{{ end }}

{{ define "course_helper_problem" }}
  {{ template "problem" loadProblem . }}
{{ end }}
