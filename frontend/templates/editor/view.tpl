{{ define "editor_view" }}
<section class="editor">
  <article>
    <div class="editor-files">
      <input type="submit" value="Ny fil +" class="btn-green outlined" onclick="newFile();">
    </div>
    <div class="editor-main">
      <div id="editfield" style="width:100%;flex: 1 0 auto;border:1px solid grey; margin-bottom: 15px;"></div>
    </div>
    <div class="editor-running">
      <input type="submit" value="Kör programmet" class="btn-green outlined">
      <div class="form-group">
        <span class="select-field" style="width: auto; margin-right: 10px;">
          <select id="language-selector" name="language" onchange="updateEditor();">
            {{ range .D.Languages }}
            <option value="{{ .LanguageId }}" data-lang="{{ .VsName }}">{{ .Name }} ({{ .Version }})</option>
            {{ end }}
          </select>
        </span>
      </div>
      <textarea name="input" style="width: 100%; flex: 1"></textarea>
      <textarea name="output" style="width: 100%; flex: 1"></textarea>
    </div>
  </article>
</section>
<script src="/static/vs/loader.js"></script>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
<script>
    var editor;
    require.config({ paths: { 'vs': '/static/vs' }});
    require(['vs/editor/editor.main'], function() {
        editor = monaco.editor.create(document.getElementById('editfield'), {
            value: [
            ].join('\n'),
        });
        updateEditor();
    });
    function updateEditor() {
			var model = editor.getModel(); // we'll create a model for you if the editor created from string value.
			var lang = document.getElementById("language-selector").selectedOptions[0].dataset.lang
			monaco.editor.setModelLanguage(model, lang);
    }

    function updateFiles(data) {
    }

    function refreshFiles() {
      $.get("/api/editor/files", function(data){
        updateFiles(data);
      }).fail(function() {
        alert("Kunde inte lista dina filer");
      });
    }

    function newFile() {
      var fileName = prompt("Vad ska filen heta?");
      $.post("/api/editor/files", {name: fileName}, function(data){
        refreshFiles();
      }).fail(function() {
        alert("Kunde inte skapa en fil");
      });
    }

</script>
{{ end }}
