{{ define "submission_status" }}
    {{ $status := .run.StatusString .version .filtered }}
    {{ if .run.Accepted .version }}
        <span class="text-col-green">
            <strong>{{ $status }}</strong>
        </span>
    {{ end }}
    {{ if .run.PartialAccepted .version }}
        <div class="text-col-yellow">
            <strong>{{ $status }}</strong>
        </div>
    {{ end }}
    {{ if .run.Rejected }}
        <span class="text-col-red">
            <strong>{{ $status }}</strong>
        </span>
    {{ end }}
    {{ if .run.Waiting }}
        <span class="submission-waiting">
            <i class="material-icons">timer</i>
            <strong>{{ $status }}</strong>
        </span>
    {{ end }}
{{ end }}

{{ define "submission_list" }}
    <table class="mdl-data-table mdl-js-data-table mdl-data-table--selectable mdl-shadow--2dp" style="width: 100%;">
        <thead>
        <tr>
            <th>ID</th>
            <th class="mdl-data-table__cell--non-numeric">Problem</th>
            <th class="mdl-data-table__cell--non-numeric">Inskickningstid</th>
            {{ if not .filtered }}
                <th class="mdl-data-table__cell--non-numeric">Språk</th>
            {{ end }}
            <th class="mdl-data-table__cell--non-numeric">Resultat</th>
        </tr>
        </thead>
        {{range .submissions }}
            {{$prob := index $.problems .ProblemID }}
            <tr>
                <td>
                    {{ if not $.filtered }}
                        <a href="{{ .Link }}">{{ .SubmissionID }}</a>
                    {{ else }}
                        {{ .SubmissionID }}
                    {{ end }}
                </td>
                <td class="mdl-data-table__cell--non-numeric">
                    <a href="{{ $prob.Link }}">{{ $prob.LocalizedTitle $.C.Locales }}</a>
                </td>
                <td class="mdl-data-table__cell--non-numeric">
                    {{ .Created.Format "2006-01-02 15:04:05"  }}
                </td>
                {{ if not $.filtered }}
                    <td class="mdl-data-table__cell--non-numeric">{{ (language .Language).Name }}</td>
                {{ end }}
                <td class="mdl-data-table__cell--non-numeric">
                    {{ template "submission_status" dict "run" .CurrentRun "filtered" $.filtered "version" $prob.CurrentVersion}}
                </td>
            </tr>
        {{end}}
    </table>
{{ end }}
gt