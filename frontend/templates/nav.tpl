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
        <li class="active"><a href="/problems">Problem</a></li>
        <li class="active"><a href="/login">Logga in</a></li>
        <li class="active"><a href="/register">Skapa konto</a></li>
      </ul>
    </nav>
  </div>
</header>
{{ end }}
