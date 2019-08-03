{{ define "course_chapter_box" }}
<a href="{{ .Chapter.Link }}">
  <div class="chapter-box">
    <h2>{{ (.Chapter.Loc $.C.Locales).Name }}</h2>
    <p>{{ (.Chapter.Loc $.C.Locales).Summary }}</p>
  </div>
</a>
{{ end }}

{{ define "course_next_section" }}
{{ .NextSection }}
{{ with .NextSection }}
{{ template "course_chapter_box" dict "Chapter" . "C" $.C }}
{{ end }}
{{ end }}
