{{ define "problems_view" }}
  <section class="problem">
    <article>
      {{ template "helper_contest_banner" . }}
      <div class="wide-row mdl-grid">
        <div class="mdl-cell mdl-cell--3-col">
          <div class="mdl-card mdl-shadow--2dp" style="width: 100%; height: auto; min-height: 0">
            <div class="mdl-color-text--grey-600 mdl-card__supporting-text">
              <table>
                <tr>
                  <td><strong>Maxpoäng:</strong></td>
                  <td>{{ .D.Problem.CurrentVersion.MaxScore }}</td>
                </tr>
                <tr>
                <tr>
                  <td><strong>Tidsgräns:</strong></td>
                  <td>{{ .D.Problem.CurrentVersion.TimeLimString }}</td>
                </tr>
                <tr>
                  <td><strong>Minnesgräns:</strong></td>
                  <td>{{ .D.Problem.CurrentVersion.MemLimString }}</td>
                </tr>
              </table>
            </div>
            <strong>
            </strong>
            <div class="mdl-card__actions mdl-card--border">
              {{ if not .C.User }}
                <a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="{{ path "login" }}">
                  Logga in för att skicka in lösningar
                </a>
              {{ else if and .C.Contest (not .C.Team) }}
                  {{ template "helper_contest_register" "Anmäl dig för att skicka in lösningar"}}
                {{ else }}
                  <a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="{{ .D.Problem.SubmitLink }}">
                    Skicka in
                  </a>
                {{ end }}
            </div>
          </div>

        </div>
        <div class="mdl-shadow--2dp mdl-cell mdl-cell--7-col" style="padding: 0 25px">
          <h1 class="display">{{ .D.Problem.LocalizedTitle $.C.Locales }}</h1>
          <div class="problem-info">
          </div>
          <div class="statement">
            {{ .D.Problem.LocalizedStatement $.C.Locales }}

            {{ range $i, $s := .D.Problem.CurrentVersion.Samples }}
              {{ template  "problem_sample" $s }}
            {{end}}

          </div>
          <p class="problem-authors">Författare: {{ .D.Problem.Author }} | Licens: {{ .D.Problem.License }}
        </div>
      </div>
    </article>
  </section>
{{ end }}
