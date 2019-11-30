{{ define "helper_contest_banner" }}
    <header class="article-header">
        <div class="row">
            <h1 class="display">
                {{ if .C.Contest.Flexible }}
                    {{ if .C.Contest.FullOver }}
                        Tävlingen är helt avslutad
                    {{ else if not .C.Contest.FullStart }}
                        Tävlingen börjar om {{ .C.Contest.UntilStart | interval }}
                    {{ else if not (.C.Contest.Started .C.Team) }}
                        Tävlingen slutar om {{ .C.Contest.UntilFullEnd | interval }}
                    {{ else if .C.Contest.Over .C.Team }}
                        Din tävling är avslutad. Tävlingen slutar om {{ .C.Contest.UntilFullEnd | interval }}
                    {{ else }}
                        Din tävling slutar om {{ .C.Contest.UntilEnd .C.Team | interval }}
                    {{ end }}
                {{ else }}
                    {{ if .C.Contest.FullOver }}
                        <h1 class="display">Tävlingen är avslutad</h1>
                    {{ else if .C.Contest.FullStart }}
                        <h1 class="display">Tävlingen slutar om {{ .C.Contest.UntilFullEnd | interval }}</h1>
                    {{ else  }}
                        <h1 class="display">Tävlingen börjar om {{ .C.Contest.UntilFullStart | interval }}</h1>
                    {{ end }}
                {{ end }}
            </h1>

        </div>
    </header>
{{ end }}

{{ define "helper_contest_register" }}
    <form method="post" action="{{ path "contest_team_register" }}" name="register">
        <div class="form-group">
            <div class="submit-field">
                <a href="#" onclick="document.forms['register'].submit()" class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">
                    {{ . }}
                </a>
            </div>
        </div>
    </form>
{{ end }}

{{ define "helper_contest_start" }}
    <form method="post" action="{{ path "contest_team_start" }}" name="start">
        <div class="form-group">
            <div class="submit-field">
                <a href="#" onclick="document.forms['start'].submit()" class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">
                    {{ . }}
                </a>
            </div>
        </div>
    </form>
{{ end }}
