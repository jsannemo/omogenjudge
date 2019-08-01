{{ define "submissions_view" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Inskickning {{ .D.Submission.SubmissionId}} - {{ .D.Problem.LocalizedTitle $.C.Locales }}</h1>
        {{ .D.Submission.StatusString }}
      </div>
    </header>
    <div class="row">
        {{ range .D.Submission.Files }}
          <pre>{{ .Contents }}</pre>
        {{end}}
    </div>
  </article>
</section>
{{ end }}
