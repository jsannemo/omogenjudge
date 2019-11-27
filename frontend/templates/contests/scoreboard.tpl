{{ define "contest_scoreboard" }}
    <section>
        <article>
            {{ template "helper_contest_banner" .C.Contest }}
            <div class="row">
                <table style="margin: auto" class="mdl-data-table mdl-js-data-table mdl-data-table--selectable">
                    <thead>
                    <tr>
                        <th style="width: 60px">#</th>
                        <th style="width: 220px" class="mdl-data-table__cell--non-numeric">Namn</th>
                        <th style="width: 60px">
                            {{.D.MaxScore}}
                            <br>
                            Totalpoäng
                        </th>
                        {{ range .D.Problems }}
                            <th style="width: 60px; text-align: center">
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
                                <td style="text-align: center">
                                    {{ $subs := index $t.Submissions .Problem.ProblemID }}
                                    {{ if $subs }}
                                        <div>
                                            {{ index $t.Scores .Problem.ProblemID }}
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
