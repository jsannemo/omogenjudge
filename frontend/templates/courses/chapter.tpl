{{ define "courses_chapter" }}
<section class="course">
  <article>
    <header class="article-header">
      <div class="display subtext">
        <a href="{{ .D.Chapter.Course.Link }}">
          {{ (.D.Chapter.Course.Loc $.C.Locales).Name.HTML }}
        </a>
      </div>
      <h1 class="display">{{ (.D.Chapter.Loc $.C.Locales).Name.HTML }}</h1>
    </header>
    <div class="row">
      <div class="statement">
      {{ (.D.Chapter.Loc $.C.Locales).Desc.HTML }}
      </div>
      {{ range .D.Chapter.Sections }}
        {{ template "course_chapter_box" dict "Chapter" . "C" $.C }}
      {{ end }}
      </div>
  </article>
</section>
{{ end }}
