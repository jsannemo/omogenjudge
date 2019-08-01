{{ define "problems_list" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Problem</h1>
      </div>
    </header>
    <div class="row">
      <table>
        <tr>
          <th>ID</th>
          <th>Title</th>
        </tr>
			 {{range .D.Problems}}
        <tr>
          <td>{{ .ShortName }}</td>
          <td>
            <a href="{{ .Link }}">{{ .LocalizedTitle $.C.Locales }}</a>
          </td>
        </tr>
			 {{end}}
      </table>
    </div>
  </article>
</section>
{{ end }}
