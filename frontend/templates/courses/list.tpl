{{ define "courses_list" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Kurser</h1>
      </div>
    </header>
    <div class="row">
      <table>
        <tr>
          <th>Namn</th>
        </tr>
			 {{range .D.Courses}}
        <tr>
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
