{{ define "users_register" }}
<section>
  <article>
    <header class="article-header">
      <div class="row">
        <h1 class="display">Skapa konto</h1>
      </div>
    </header>
    <div class="row">
      <form method="post">
        <div class="form-group">
          <label>
            <input type="text" required name="username">
            <span class="placeholder">Användarnamn</span>
          </label>

          <label>
            <input type="password" required name="password">
            <span class="placeholder">Lösenord</span>
          </label>

          <input type="submit" value="Skapa konto" class="btn-green outline">
        </div>
      </form>
    </div>
  </article>
</section>
{{ end }}
