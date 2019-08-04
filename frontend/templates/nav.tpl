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
        <li><a href="/">Hem</a></li>
        <li><a href="/courses">Kurser</a></li>
        <li><a href="/problems">Uppgiftsarkiv</a></li>
        {{ if .C.User }}
          <li><a href="/users/{{ .C.User.Username}}">{{ .C.User.Username }}</a></li>
        {{ else }}
          <li><a href="/login">Logga in</a></li>
          <li><a href="/register">Skapa konto</a></li>
        {{ end }}
      </ul>
    </nav>
  </div>
</header>
<section class="content">
{{ end }}
