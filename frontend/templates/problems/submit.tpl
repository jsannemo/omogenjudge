{{ define "problems_submit" }}
    <section>
        <article>
            {{ template "helper_contest_banner" . }}
            <div class="row">
                <h1 class="display">Skicka in - {{ .D.Problem.LocalizedTitle $.C.Locales }}</h1>
                <form method="POST" onsubmit="return onSubmit()">
                    <div id="submitfield" style="width:100%;height:550px;border:1px solid grey; margin-bottom: 15px;"></div>
                    <textarea style="display: none" id="submission" name="submission"></textarea>
                    <div class="form-group">
            <span class="select-field" style="width: auto; margin-right: 10px;">
              <select id="language-selector" name="language" onchange="updateEditor();">
                {{ range .D.Languages }}
                    <option value="{{ .LanguageId }}" data-lang="{{ .VsName }}">{{ .Name }} ({{ .Version }})</option>
                {{ end }}
              </select>
            </span>
                        <input type="submit" value="Skicka in" class="btn-green outlined">
                    </div>
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
            });
            updateEditor();
        });
        function onSubmit() {
            var code = editor.getValue();
            document.getElementById('submission').value = code;
        }
        function updateEditor() {
            var model = editor.getModel(); // we'll create a model for you if the editor created from string value.
            var lang = document.getElementById("language-selector").selectedOptions[0].dataset.lang
            console.log(lang);
            monaco.editor.setModelLanguage(model, lang);
        }
    </script>
{{ end }}
