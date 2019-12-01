{{ define "home_home" }}
    <article>
        {{ template "helper_contest_banner" . }}
        {{ if not (.C.Contest.Started .C.Team) }}
            {{ template "home_before_contest" . }}
        {{ else }}
            {{ template "home_during_contest" . }}
        {{ end }}
    </article>
{{ end }}

{{ define "home_during_contest" }}
    <div class="row mdl-grid">
        {{ range .C.Contest.Problems }}
                <div class="mdl-cell mdl-cell--6-col">
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
                                    <th style="padding-left: 0px; line-height: 14px">
                                        Maxpoäng<br>
                                        Grupp
                                    </th>
                                    {{ range $i, $g := $p.Groups}}
                                        {{ if not $g.PublicVisibility }}
                                            <th style="padding-left: 0px; line-height: 14px">
                                                {{$g.Score}}<br>
                                                G{{ $i }}
                                            </th>
                                        {{ end }}
                                    {{ end }}
                                    <th style="padding-left: 0px; line-height: 14px">
                                        {{ .Problem.CurrentVersion.MaxScore }}<br>
                                        TOT
                                    </th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr>
                                    <td>Din poäng</td>
                                    {{ range $i, $g := $p.Groups}}
                                        {{ if not $g.PublicVisibility }}
                                            <td style="padding-left: 0px">
                                                {{ $score := index $p.Scores $g.Name }}
                                                {{ $bg := ""}}
                                                {{ if eq $score $g.Score }}
                                                    {{ $bg = "text-col-green" }}
                                                {{ else if eq $score 0}}
                                                    {{ $bg = "text-col-red"}}
                                                {{ else }}
                                                    {{ $bg = "text-col-yellow"}}
                                                {{ end }}
                                                <span class="{{ $bg }}">
                                                    {{ $score }}
                                                </span>
                                            </td>
                                        {{ end }}
                                    {{ end }}
                                    <td style="padding-left: 0px">
                                        {{ $score := index $p.Score }}
                                        {{ $bg := ""}}
                                        {{ if eq $score .Problem.CurrentVersion.MaxScore }}
                                            {{ $bg = "text-col-green" }}
                                        {{ else if eq $score 0}}
                                            {{ $bg = "text-col-red"}}
                                        {{ else }}
                                            {{ $bg = "text-col-yellow"}}
                                        {{ end }}
                                        <strong class="{{ $bg }}">{{ $p.Score }}</strong>
                                    </td>
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
                            {{ template "helper_contest_register" "Anmäl dig för att skicka in lösningar"}}
                        {{ else }}
                            <a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" href="{{ .Problem.SubmitLink }}">
                                Skicka in lösning
                            </a>
                        {{ end }}
                    </div>
                </div>
            </div>
        {{ end }}
        {{ if and (not .C.Contest.FullOver) (not .C.Team) }}
            <h1>Anmälan</h1>
            {{ if not .C.User }}
                Du måste <a href="/login">logga in</a> för att anmäla dig till tävlingen.
            {{ else }}
                För att kunna delta i tävlingen måste du först anmäla dig.
                {{ template "helper_contest_register" "Anmäl dig"}}
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
                    <td>{{ .C.Contest.FullEndTime | date "2006-01-02 15:04" }}</td>
                </tr>
                <tr>
                    <th>Längd:</th>
                    <td>{{ .C.Contest.Duration | interval }}</td>
                </tr>
            </table>
        </div>
        <div style="clear: both"></div>
        {{ if not .C.Team }}
            <h1>Anmälan</h1>
            {{ if not .C.User }}
                Du måste <a href="/login">logga in</a> för att anmäla dig till tävlingen.
            {{ else if not .C.Team }}
                För att kunna delta i tävlingen måste du först anmäla dig.
                {{ template "helper_contest_register" "Anmäl dig"}}
            {{ end }}
        {{ else if .C.Contest.Flexible }}
            <h1>Starta tävlingen</h1>
            När du startar tävlingen har du {{ .C.Contest.Duration | interval }} på dig att slutföra tävlingen.
            {{ template "helper_contest_start" "Starta tävlingen"}}
        {{ else }}
            <h1>Anmälan</h1>
            Du är redan anmäld till tävlingen.
        {{ end }}
    </div>
{{ end }}
