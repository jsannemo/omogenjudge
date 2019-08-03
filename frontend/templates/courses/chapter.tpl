{{ define "courses_chapter" }}
<section class="course">
  <article>
    <header class="article-header">
      <span class="display subtext">
        <a href="{{ .D.Chapter.Course.Link }}">
          {{ (.D.Chapter.Course.Loc $.C.Locales).Name.HTML }}
        </a>
      </span>
      <h1 class="display">{{ (.D.Chapter.Loc $.C.Locales).Name.HTML }}</h1>
    </header>
    <div class="row">
      {{ (.D.Chapter.Loc $.C.Locales).Desc.HTML }}
      {{ range .D.Chapter.Sections }}
        {{ template "course_chapter_box" dict "Chapter" . "C" $.C }}
      {{ end }}
      </div>
  </article>
</section>
{{ end }}
