{{ define "home_home" }}
    <article>
        {{ template "helper_contest_banner" .C.Contest }}
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
            <div style="clear: both" action="/team/register" method="POST"></div>
            <h1>Anmälan</h1>
            {{ if not .C.User }}
                Du måste <a href="/login">logga in</a> för att anmäla dig till tävlingen.
            {{ else if not .C.Team }}
                För att kunna delta i tävlingen måste du först anmäla dig.
                <form method="post">
                    <div class="form-group">
                        <div class="submit-field">
                            <input type="submit" value="Anmäl dig" class="raised">
                        </div>
                    </div>
                </form>
            {{ else }}
                Du är redan anmäld till tävlingen.
            {{ end }}
        </div>
    </article>
{{ end }}
