{{ define "users_view" }}
    <section>
        <article>
            {{ template "helper_contest_banner" .C.Contest }}
            <div class="row">
                <h1 class="display">Inskickningar - {{ .D.Username }}</h1>
                {{ template "submission_list" dict "submissions" .D.Submissions "problems" .D.Problems "C" .C "filtered" .D.Filtered }}
            </div>
        </article>
    </section>
{{ end }}
