{{ define "courses_list" }}
<section class="course">
  <article>
    <header class="article-header">
      <h1 class="display">Kurser</h1>
    </header>
    <div class="row">
      {{ range .D.Courses }}
        {{ template "course_chapter_box" dict "Chapter" . "C" $.C }}
      {{ end }}
      </div>
    </div>
  </article>
</section>
{{ end }}
