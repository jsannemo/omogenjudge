{{ define "submission_list" }}
<table class="bordered" style="width: 100%;">
<thead>
  <tr>
    <th>ID</th>
    <th>Problem</th>
    <th>Inskickningstid</th>
    <th>Språk</th>
    <th>Resultat</th>
  </tr>
</thead>
 {{range .submissions }}
 {{$prob := index $.problems .ProblemId }}
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
{{ end }}
