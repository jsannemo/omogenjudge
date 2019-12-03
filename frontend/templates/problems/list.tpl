{{ define "problems_list" }}
  <section>
    <article>
      <div class="row mdl-griw">
        <div class="mdl-cell mdl-cell--12-col">
          <h1>Problem</h1>
        </div>
        <div class="mdl-cell mdl-cell--12-col">
          <table class="mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp">
            <thead>
            <tr>
              <th class="mdl-data-table__cell--non-numeric">ID</th>
              <th class="mdl-data-table__cell--non-numeric">Title</th>
            </tr>
            </thead>
            {{range .D.Problems}}
              <tr>
                <td class="mdl-data-table__cell--non-numeric">{{ .ShortName }}</td>
                <td class="mdl-data-table__cell--non-numeric">
                  <a href="{{ .Link }}">{{ .LocalizedTitle $.C.Locales }}</a>
                </td>
              </tr>
            {{end}}
          </table>
        </div>
    </article>
  </section>
{{ end }}
