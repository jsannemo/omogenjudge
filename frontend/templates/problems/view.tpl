{{ define "problems_view" }}
  <section class="problem">
    <article>
      {{ template "helper_contest_banner" . }}
      <div class="wide-row mdl-grid">
        <div class="mdl-cell mdl-cell--3-col">
          <div class="mdl-card mdl-shadow--2dp" style="width: 100%; height: auto; min-height: 0">
            <div class="mdl-color-text--grey-600 mdl-card__supporting-text">
              <table>
                {{ if .D.Problem.CurrentVersion.MaxScore }}
                  <tr>
                    <td><strong>Maxpoäng:</strong></td>
                    <td>{{ .D.Problem.CurrentVersion.MaxScore }}</td>
                  </tr>
                {{ end }}
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
          {{ if .D.Problem.StatementFiles }}
            <div class="mdl-card mdl-shadow--2dp" style="margin-top: 25px; padding: 10px 25px; width: auto; height: auto; min-height: 0">
                <strong>Bifogade filer</strong>
                {{ range .D.Problem.StatementFiles }}
                  <a href="{{ path "problem_file" "problem_name" $.D.Problem.ShortName "problem_file_name" .Path }}">{{ .Path }}</a>
                {{ end }}
            </div>
          {{ end }}
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
    <script>
      new ClipboardJS('.copy-sample-btn', {
        target: function(trigger) {
          showTooltip(trigger,'Copied!');
          return trigger.parentElement.nextElementSibling;
        }
      });
    </script>
  </section>
{{ end }}
