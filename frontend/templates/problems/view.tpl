{{ define "problems_view" }}
<section class="course problem">
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">{{ .D.Problem.LocalizedTitle $.C.Locales }}</h1>
        <div class="problem-info">Tidsgräns {{ .D.Problem.TimeLimString }} | Minnesgräns {{ .D.Problem.MemLimString }}</div>
        <div class="problem-actions"> <a class="button" href="{{ .D.Problem.SubmitLink }}">Skicka in</a></div>
      </div>
    </header>
    <div class="row">
      <div class="statement">
        {{ .D.Problem.LocalizedStatement $.C.Locales }}

        {{ range $i, $s := .D.Problem.Samples }}
          {{ template  "problem_sample" $s }}
        {{end}}

      </div>
      <p class="problem-authors">Författare: {{ .D.Problem.Author }} | Licens: {{ .D.Problem.License }}
    </div>
  </article>
</section>
{{ end }}
