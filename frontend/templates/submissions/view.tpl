{{ define "submissions_view" }}
  <article>
    {{ template "helper_contest_banner" .C.Contest }}
    <div class="row">
      {{ template "submission_list" dict "submissions" (list .D.Submission) "problems" .D.Problems "C" .C "filtered" false }}
      <table class="bordered" style="width: 100%">
        <thead>
        <tr>
          <th> Testgrupp </th>
          <th> Resultat </th>
          <th> Poäng </th>
        </tr>
        </thead>
        {{ range $_, $p := .D.Problems }}
        {{ range $i, $g := $p.CurrentVersion.TestGroups }}
            <tr>
                <td>
                  {{ if $g.PublicVisibility }}
                    <strong>Exempelfall</strong>
                  {{ else }}
                    <strong>Grupp {{ $i }}</strong>
                  {{ end }}
                </td>
                <td>
                  {{ $.D.Submission.CurrentRun.GroupVerdict $g.TestGroupID }}
                </td>
                <td>
                  {{ $.D.Submission.CurrentRun.GroupScore $g.TestGroupID }}
                </td>
            </tr>
        {{ end }}
      {{ end }}
      </table>
      {{ if .D.Submission.CurrentRun.CompileError.Valid }}
        <table class="bordered" style="width: 100%; margin-top: 15px;">
          <thead>
          <tr><th>Felmeddelanden från kompilatorn</th></tr>
          </thead>
          <tr><td><pre>{{ .D.Submission.CurrentRun.CompileError.Value }}</pre></td></tr>
        </table>
      {{ end }}
      {{ range .D.Submission.Files }}
        <table class="bordered" style="width: 100%; margin-top: 15px;">
          <thead>
          <tr><th>{{ .Path }}</th></tr>
          </thead>
          <tr><td><pre><code data-lang="{{ (language $.D.Submission.Language).VsName }}" class="code-colorize">{{ .Contents }}</code></pre></td></tr>
        </table>
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
