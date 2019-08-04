{{ define "courses_section" }}
<section class="course">
  <article>
    <header class="article-header">
      <div class="display subtext">
        <a href="{{ .D.Section.Chapter.Course.Link }}">
          {{ (.D.Section.Chapter.Course.Loc $.C.Locales).Name }}
        </a>
        /
        <a href="{{ .D.Section.Chapter.Link }}">
            {{ (.D.Section.Chapter.Loc $.C.Locales).Name }}
        </a>
      </div>
      <h1 class="display">{{ (.D.Section.Loc $.C.Locales).Name }}</h1>
    </header>
    <div class="row">
    <div class="statement">
    {{ .D.Output }}
    </div>

    {{ template "course_next_section" dict "section" .D.Section "C" $.C }}
    </div>
  </article>
</section>
{{ end }}
