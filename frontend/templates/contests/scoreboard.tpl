{{ define "contest_scoreboard" }}
    <section>
        <article>
            {{ template "helper_contest_banner" . }}
            <div class="row">
                <table style="width: 100%; margin: auto" class="mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp">
                    <thead>
                    <tr>
                        <th style="width: 60px">#</th>
                        <th style="width: auto" class="mdl-data-table__cell--non-numeric">Namn</th>
                        <th style="width: 80px">
                            {{.D.MaxScore}}
                            <br>
                            Totalpoäng
                        </th>
                        {{ range .D.Problems }}
                            <th style="width: 80px; text-align: center">
                                {{.Problem.CurrentVersion.MaxScore}}
                                <br>
                                <a href="{{ .Problem.Link }}">
                                    {{ .Label }}
                                </a>
                            </th>
                        {{ end }}
                    </tr>
                    </thead>
                    <tbody>
                    {{ range $_, $t := .D.Teams }}
                        <tr>
                            <td>{{ .Rank }}</td>
                            <td class="mdl-data-table__cell--non-numeric">
                                <a href="{{ .Team.Link }}">
                                    {{ .Team.DisplayName }}
                                </a>
                            </td>
                            <td>
                                <strong>{{ $t.TotalScore }}</strong>
                            </td>
                            {{ range $.D.Problems }}
                                {{ $subs := index $t.Submissions .Problem.ProblemID }}
                                {{ $score := index $t.Scores .Problem.ProblemID }}
                                {{ $scoreCol := "white"}}
                                {{ if $subs }}
                                    {{ $scoreCol = index $t.ScoreCols .Problem.ProblemID }}
                                {{ end }}

                                <td style="background-color: {{ $scoreCol }}; text-align: center; border: 1px solid rgba(0,0,0,.12); padding: 0">
                                    {{ if $subs }}
                                        <div>
                                            {{ $score }}
                                        </div>
                                        <div style="font-size: 12px">
                                            {{ index $t.Times .Problem.ProblemID | hhmm }}
                                        </div>
                                    {{ else }}
                                    {{ end }}
                                </td>
                            {{ end }}
                        </tr>
                    {{ end }}
                    </tbody>
                </table>
            </div>
        </article>
    </section>
{{ end }}
