{{ define "nav" }}
<header class="navbar">
  <div class="navbar-container">
    <div class="navbar-top">
      <img class="navbar-logo" src="/static/kodsport/logo.svg">
      <button class="navbar-hamburger">
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
    </div>
    <nav class="navbar-nav closed">
      <ul>
        <li class="active"><a href="/">Hem</a></li>
        <li class="active"><a href="/courses">Kurser</a></li>
        <li class="active"><a href="/problems">Uppgiftsarkiv</a></li>
        {{ if .C.User }}
          <li class="active"><a href="/users/{{ .C.User.Username}}">{{ .C.User.Username }}</a></li>
        {{ else }}
          <li class="active"><a href="/login">Logga in</a></li>
          <li class="active"><a href="/register">Skapa konto</a></li>
        {{ end }}
      </ul>
    </nav>
  </div>
</header>
{{ end }}
