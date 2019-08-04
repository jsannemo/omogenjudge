{{ define "nav" }}
<header class="navbar">
  <div class="navbar-container">
    <div class="navbar-top">
    <a href="/"><img class="navbar-logo" src="/static/kodsport/logo.svg"></a>
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
        <li><a href="/problems">Problemarkiv</a></li>
        {{ if .C.User }}
          <li class="navbar-dropdown closed">
            <a href="javascript:;"> <i class="material-icons">person</i> {{ .C.User.Username }}<span class="navbar-dropdown-caret"></span></a>
            <ul>
              <li><a href="/users/{{ .C.User.Username}}">Inskickningar</a>
              <li><a href="/logout">Logga ut</a></li>
            <ul>
          </li>
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
