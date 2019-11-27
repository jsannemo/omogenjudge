{{ define "home_home" }}
    <article>
        {{ template "helper_contest_banner" .C.Contest }}
        {{ if not .C.Contest.Started}}
            {{ template "home_before_contest" . }}
        {{ else }}
            {{ template "home_during_contest" . }}
        {{ end }}
    </article>
{{ end }}

{{ define "home_during_contest" }}
    <div class="row">
        {{ range .C.Contest.Problems }}
            <div class="mdl-grid">
                <div class="mdl-card mdl-shadow--2dp" style="width: 100%; height: auto; min-height: 0">
                    <div class="mdl-card__title">
                        <h2 class="mdl-card__title-text"><strong>{{ .Label }}</strong>&nbsp;{{ .Problem.LocalizedTitle $.C.Locales }}</h2>
                    </div>
                    {{ if $.D }}
                        {{ $p := index $.D.Problems .ProblemID }}
                        <div class="mdl-card__supporting-text" style="width: 100%">
                            <table style="margin: auto" class="mdl-data-table mdl-js-data-table mdl-data-table--selectable">
                                <thead>
                                <tr>
                                    <td></td>
                                    {{ range $i, $g := $p.Groups}}
                                        {{ if not $g.PublicVisibility }}
                                            <th>
                                                G{{ $i }}
                                                ({{$g.Score}})
                                            </th>
                                        {{ end }}
                                    {{ end }}
                                    <th>TOT ({{ .Problem.CurrentVersion.MaxScore }})</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr>
                                    <th>Din poäng</th>
                                    {{ range $i, $g := $p.Groups}}
                                        {{ if not $g.PublicVisibility }}
                                            <td>
                                                {{ index $p.Scores $g.Name }}
                                            </td>
                                        {{ end }}
                                    {{ end }}
                                    <td><strong>{{ $p.Score }}</strong></td>
                                </tr>
                                </tbody>
                            </table>

                        </div>
                    {{ end }}
                    <div class="mdl-card__actions mdl-card--border">
                        <a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="{{ .Problem.Link }}">
                            Öppna problemlydelsen
                        </a>
                        {{ if not $.C.Team }}
                            <a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="{{ path "contest_team_register" }}">
                                Anmäl dig för att skicka in lösningar
                            </a>
                        {{ else }}
                            <a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="{{ .Problem.SubmitLink }}">
                                Skicka in lösning
                            </a>
                        {{ end }}
                    </div>
                </div>
            </div>
        {{ end }}
        {{ if and (not .C.Contest.Over) (not .C.Team) }}
            <h1>Anmälan</h1>
            {{ if not .C.User }}
                Du måste <a href="/login">logga in</a> för att anmäla dig till tävlingen.
            {{ else }}
                För att kunna delta i tävlingen måste du först anmäla dig.
                <form method="post" action="{{ path "contest_team_register" }}" method="POST">
                    <div class="form-group">
                        <div class="submit-field">
                            <input type="submit" value="Anmäl dig" class="raised">
                        </div>
                    </div>
                </form>
            {{ end }}
        {{ end }}
    </div>
{{ end }}

{{ define "home_before_contest" }}
    <div class="row">
        <h1>{{ .C.Contest.Title }}</h1>
        <div>
            <span>
                    Välkommen till tävlingssystemet för <strong>{{ .C.Contest.Title }}</strong>.
                </span>
            <table class="bordered" style="float: right">
                <tr>
                    <th>Starttid:</th>
                    <td>{{ .C.Contest.StartTime.Time | date "2006-01-02 15:04" }}</td>
                </tr>
                <tr>
                    <th>Sluttid:</th>
                    <td>{{ .C.Contest.EndTime | date "2006-01-02 15:04" }}</td>
                </tr>
                <tr>
                    <th>Längd:</th>
                    <td>{{ .C.Contest.Duration | interval }}</td>
                </tr>
            </table>
        </div>
        <div style="clear: both"></div>
        {{ if or (not .C.Contest.Started) (not .C.Team) }}
            <h1>Anmälan</h1>
            {{ if not .C.User }}
                Du måste <a href="/login">logga in</a> för att anmäla dig till tävlingen.
            {{ else if not .C.Team }}
                För att kunna delta i tävlingen måste du först anmäla dig.
                <form method="post" action="{{ path "contest_team_register" }}" method="POST">
                    <div class="form-group">
                        <div class="submit-field">
                            <input type="submit" value="Anmäl dig" class="raised">
                        </div>
                    </div>
                </form>
            {{ else }}
                Du är redan anmäld till tävlingen.
            {{ end }}
        {{ end }}
    </div>
{{ end }}
