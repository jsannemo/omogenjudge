{{ define "submissions_view" }}
  <article>
    {{ template "helper_contest_banner" . }}
    <div class="mdl-grid row">
        <h1>Inskickning {{ .D.Submission.SubmissionID }}</h1>
      {{ template "submission_list" dict "submissions" (list .D.Submission) "problems" .D.Problems "C" .C "filtered" .D.Filtered "queue" false }}
      {{ if .D.Submission.CurrentRun.CompileError.Valid }}
        {{ if .D.WithCode }}
          <table class="mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp" style="width: 100%; margin-top: 15px">
            <thead>
            <tr><th class="mdl-data-table__cell--non-numeric">Felmeddelanden från kompilatorn</th></tr>
            </thead>
            <tr><td class="mdl-data-table__cell--non-numeric"><pre>{{ .D.Submission.CurrentRun.CompileError.Value }}</pre></td></tr>
          </table>
        {{ end}}
      {{ else }}
        <table class="mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp" style="width: 100%; margin-top: 15px">
          <thead>
          <tr>
            <th class="mdl-data-table__cell--non-numeric"> Testgrupp </th>
            <th> Poäng </th>
          </tr>
          </thead>
          {{ range $_, $p := .D.Problems }}
            {{ range $i, $g := $p.CurrentVersion.TestGroups }}
              {{ $score := $.D.Submission.CurrentRun.GroupScore $g.Name }}
              {{ $verdict := $.D.Submission.CurrentRun.GroupVerdict $g.Name }}
              {{ $bg := "white "}}
              {{ if not $verdict.Waiting}}
                {{ if and (eq $score $g.Score) $verdict.Accepted }}
                  {{ $bg = "bg-green-100" }}
                {{ else if eq $score 0}}
                  {{ $bg = "bg-red-100"}}
                {{ else }}
                  {{ $bg = "bg-yellow-50"}}
                {{ end }}
              {{ end }}
              <tr class="{{ $bg }}">
                <td class="mdl-data-table__cell--non-numeric">
                  {{ if $g.PublicVisibility }}
                    Exempelfall
                  {{ else if not (eq $g.Score 0) }}
                    Grupp {{ $i }}
                  {{ else }}
                    Testdata
                  {{ end }}
                </td>
                <td>
                  {{ if and $verdict.Accepted (not (eq $score 0)) }}
                    {{ $score }}
                  {{ else }}
                    {{ $verdict.Filtered $.D.Filtered }}
                  {{ end }}
                </td>
              </tr>
            {{ end }}
          {{ end }}
        </table>
      {{ end }}
      {{ if .D.WithCode }}
        {{ range .D.Submission.Files }}
          <table class="mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp" style="width: 100%; margin-top: 15px;">
            <thead>
            <tr><th class="mdl-data-table__cell--non-numeric">{{ .Path }}</th></tr>
            </thead>
            <tr><td class="mdl-data-table__cell--non-numeric"><pre><code data-lang="{{ (language $.D.Submission.Language).VsName }}" class="code-colorize">{{ .Contents }}</code></pre></td></tr>
          </table>
        {{ end }}
      {{ end }}
    </div>
  </article>

  <script src="/static/vs/loader.js"></script>
  <script>
    require.config({ paths: { 'vs': '/static/vs' }});
    require(['vs/editor/editor.main'], function() {
      Array.from(document.getElementsByClassName('code-colorize')).forEach(
              d => monaco.editor.colorizeElement(d));
    });
  </script>
{{ end }}
