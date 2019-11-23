{{ define "contests_team_list" }}
  <section>
    <article>
      {{ template "helper_contest_banner" .C.Contest }}
      <div class="row">
        <table class="bordered" style="margin: auto; width: 100%">
          <thead>
          <tr>
            <th>Lag</th>
          </tr>
          </thead>
          {{range .D.Teams}}
            <tr>
              <td>{{ .DisplayName }}</td>
            </tr>
          {{end}}
        </table>
      </div>
    </article>
  </section>
{{ end }}
