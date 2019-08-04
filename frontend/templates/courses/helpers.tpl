{{ define "course_chapter_box" }}
<a href="{{ .Chapter.Link }}">
  <div class="chapter-box">
    <h2>{{ (.Chapter.Loc $.C.Locales).Name.HTML }}</h2>
    <p>{{ (.Chapter.Loc $.C.Locales).Summary.HTML }}</p>
  </div>
</a>
{{ end }}

# Nästa avsnitt
{{ define "course_next_section" }}
{{ with .section.NextSection }}
<h2>Next section</h2>
{{ template "course_chapter_box" dict "Chapter" . "C" $.C }}
{{ end }}
{{ end }}
