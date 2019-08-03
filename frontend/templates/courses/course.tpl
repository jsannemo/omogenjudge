{{ define "courses_course" }}
<section class="course">
  <article>
    <header class="article-header">
      <h1 class="display">{{ (.D.Course.Loc $.C.Locales).Name }}</h1>
    </header>
    <div class="row">
      {{ (.D.Course.Loc $.C.Locales).Desc.HTML }}
      {{ range .D.Course.Chapters }}
        {{ template "course_chapter_box" dict "Chapter" . "C" $.C }}
      {{ end }}
      </div>
  </article>
</section>
{{ end }}
