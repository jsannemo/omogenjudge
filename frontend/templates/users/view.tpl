{{ define "users_view" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Inskickningar</h1>
      </div>
    </header>
    <div class="row">
      {{ template "submission_list" dict "submissions" .D.Submissions "problems" .D.Problems "C" .C }}
    </div>
  </article>
</section>
{{ end }}
