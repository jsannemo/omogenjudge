{{ define "submission_status" }}
    {{ if .run.Accepted }}
        <span class="text-col-green">
            <i class="material-icons">done</i>
            <strong>{{ .run.StatusString .filtered }}</strong>
        </span>
    {{ end }}
    {{ if .run.Rejected }}
        <span class="text-col-red">
            <i class="material-icons">warning</i>
            <strong>{{ .run.StatusString .filtered }}</strong>
        </span>
    {{ end }}
    {{ if .run.Waiting }}
        <span class="submission-waiting">
            <i class="material-icons">timer</i>
            <strong>{{ .run.StatusString .filtered }}</strong>
        </span>
    {{ end }}
{{ end }}

{{ define "submission_list" }}
    <table class="bordered" style="width: 100%;">
        <thead>
        <tr>
            <th>ID</th>
            <th>Problem</th>
            <th>Inskickningstid</th>
            {{ if not .filtered }}
                <th>Språk</th>
            {{ end }}
            <th>Resultat</th>
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
                <td>
                    <a href="{{ $prob.Link }}">{{ $prob.LocalizedTitle $.C.Locales }}</a>
                </td>
                <td>
                    {{ .Created.Format "2006-01-02 15:04:05"  }}
                </td>
                {{ if not $.filtered }}
                    <td>{{ (language .Language).Name }}</td>
                {{ end }}
                <td align="center">
                    {{ template "submission_status" dict "run" .CurrentRun "filtered" $.filtered}}
                </td>
            </tr>
        {{end}}
    </table>
{{ end }}
gt