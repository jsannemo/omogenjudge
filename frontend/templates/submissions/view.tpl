{{ define "submissions_view" }}
<article>
  <header class="article-header">
    <div class="row">
      <h1 class="display">Inskickning {{ .D.Submission.SubmissionId }}</h1>
    </div>
  </header>
  <div class="row">
    {{ template "submission_list" dict "submissions" (list .D.Submission) "problems" .D.Problems "C" .C }}
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
