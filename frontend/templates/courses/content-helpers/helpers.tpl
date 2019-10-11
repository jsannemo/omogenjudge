{{ define "problem" }}
<div class="course-problem">
  <a href="{{ .Link }}">
  <h1 class="course-problem-title">
{{ .LocalizedTitle ctx.Locales }}
</h1>
</a>
  {{ .LocalizedStatement ctx.Locales }}

  {{ range .Samples }}
    <b>Exempel</b>
    <pre>{{ .InputFile.FileString }}</pre>
    <hr>
    <pre>{{ .OutputFile.FileString }}</pre>
  {{end}}
</div>
{{ end }}

{{ define "course_helper_problem" }}
  {{ template "problem" loadProblem . }}
{{ end }}

{{ define "problemexc" }}
<div class="course-problem-exercise">
  <p>
  <strong>Problem Exercise</strong>: <a href="{{ .Link }}">{{ .LocalizedTitle ctx.Locales }} </a>
  (not solved yet)
  </p>
</div>
{{ end }}

{{ define "course_helper_pexercise" }}
  {{ template "problemexc" loadProblem . }}
{{ end }}
