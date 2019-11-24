{{ define "nav" }}
<header class="navbar">
  <div class="navbar-container">
    <div class="navbar-top">
      <a href="/">
        {{ if .C.Contest }}
          <span class="navbar-textlogo">
          {{ .C.Contest.Title }}
          </span>
        {{ else }}
          <img class="navbar-logo" src="/static/kodsport/logo.png">
        {{ end }}
      </a>
      <button class="navbar-hamburger">
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
    </div>
    <nav class="navbar-nav closed">
      <ul>
        <li><a href="/">Hem</a></li>
        {{ if .C.Contest.Started }}
          <li><a href="/scoreboard">Poängställning</a></li>
        {{ else }}
          <li><a href="/teams">Lag</a></li>
        {{ end}}
        {{ if .C.User }}
          <li class="navbar-dropdown closed">
            <a href="javascript:;"> <i class="material-icons">person</i> {{ .C.User.Username }}<span class="navbar-dropdown-caret"></span></a>
            <ul>
              {{ if .C.Contest.Started }}
                <li><a href="/users/{{ .C.User.Username}}">Inskickningar</a></li>
              {{ end }}
              <li><a href="/logout">Logga&nbsp;ut</a></li>
            </ul>
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
