{{ define "home_home" }}
    <section>
        <article>
            {{ template "helper_contest_banner" . }}
            <div class="row">
                <h1>Domarkö</h1>
                {{ template "submission_list" dict "submissions" .D.Submissions "problems" .D.Problems "C" .C "queue" true "filtered" false }}
            </div>
        </article>
    </section>
{{ end }}
