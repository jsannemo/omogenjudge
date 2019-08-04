{{ define "problems_view" }}
<section class="course">
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">{{ .D.Problem.LocalizedTitle $.C.Locales }}</h1>
        <a href="{{ .D.Problem.SubmitLink }}">Skicka in</a>
      </div>
    </header>
    <div class="row">
        {{ .D.Problem.LocalizedStatement $.C.Locales }}

        {{ range .D.Problem.Samples }}
          <b>Sample</b>
          <pre>{{ .InputFile.FileString }}</pre>
          <hr>
          <pre>{{ .OutputFile.FileString }}</pre>
        {{end}}
    </div>
  </article>
</section>
{{ end }}
