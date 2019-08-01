{{ define "problems_submit" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Skicka in - {{ .D.Problem.LocalizedTitle $.C.Locales }}</h1>
      </div>
    </header>
    <div class="row">
      <form method="POST" onsubmit="return onSubmit()">
				<div id="submitfield" style="width:100%;height:600px;border:1px solid grey"></div>
				<textarea style="display: none" id="submission" name="submission"></textarea>
        <input type="submit" value="Skicka in">
      </form>
    </div>
  </article>
</section>
<script src="/static/vs/loader.js"></script>
<script>
    var editor;
    require.config({ paths: { 'vs': '/static/vs' }});
    require(['vs/editor/editor.main'], function() {
        editor = monaco.editor.create(document.getElementById('submitfield'), {
            value: [
            ].join('\n'),
            language: 'python'
        });
    });
		function onSubmit() {
      var code = editor.getValue();
      console.log(code);
      document.getElementById('submission').value = code;
    }
</script>
{{ end }}
