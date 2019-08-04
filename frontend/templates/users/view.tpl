{{ define "users_view" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Inskickningar</h1>
      </div>
    </header>
    <div class="row">
      <table class="bordered" style="margin: auto">
      <thead>
        <tr>
          <th>ID</th>
          <th>Problem</th>
          <th>Inskickningstid</th>
          <th>Språk</th>
          <th>Resultat</th>
        </tr>
      </thead>
			 {{range .D.Submissions}}
       {{$prob := index $.D.Problems .ProblemId }}
        <tr>
          <td>
            <a href="{{ .Link }}">{{ .SubmissionId }}</a>
          </td>
          <td>
            <a href="{{ $prob.Link }}">{{ $prob.LocalizedTitle $.C.Locales }}</a>
          </td>
          <td>
            {{ .Created.Format "2006-01-02 15:04:05"  }}
          </td>
          <td>{{ (language .Language).Name }}</td>
          <td>
            {{ .StatusString }}
          </td>
        </tr>
			 {{end}}
      </table>
    </div>
  </article>
</section>
{{ end }}
