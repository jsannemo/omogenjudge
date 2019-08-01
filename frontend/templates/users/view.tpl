{{ define "users_view" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Inskickningar</h1>
      </div>
    </header>
    <div class="row">
      <table>
        <tr>
          <th>ID</th>
          <th>Problem</th>
          <th>Resultat</th>
        </tr>
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
            {{ .StatusString }}
          </td>
        </tr>
			 {{end}}
      </table>
    </div>
  </article>
</section>
{{ end }}
