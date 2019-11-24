{{ define "contest_scoreboard" }}
    <section>
        <article>
            {{ template "helper_contest_banner" .C.Contest }}
            <div class="row">
                <table style="margin: auto" class="mdl-data-table mdl-js-data-table mdl-data-table--selectable">
                    <thead>
                    <tr>
                        <th style="width: 60px">#</th>
                        <th style="width: 220px" class="mdl-data-table__cell--non-numeric">Name</th>
                        <th style="width: 60px">Totalpoäng</th>
                        {{ range .D.Problems }}
                            <th style="width: 60px">{{ .Label }}</th>
                        {{ end }}
                    </tr>
                    </thead>
                    <tbody>
                    {{ range $_, $t := .D.Teams }}
                        <tr>
                            <td>{{ .Rank }}</td>
                            <td class="mdl-data-table__cell--non-numeric">{{ .Team.DisplayName }}</td>
                            <td>
                                <strong>{{ $t.TotalScore }}</strong>
                            </td>
                            {{ range $.D.Problems }}
                                <td style="text-align: center">
                                    {{ $subs := index $t.Submissions .ProblemID }}
                                    {{ if $subs }}
                                        <div>
                                            {{ index $t.Scores .ProblemID }}
                                        </div>
                                        <div style="font-size: 12px">
                                                {{ index $t.Times .ProblemID | hhmm }}
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
