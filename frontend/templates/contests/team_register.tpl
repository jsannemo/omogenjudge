{{ define "contest_team_register" }}
	<section>
		<article>
			{{ template "helper_contest_banner" .C.Contest }}
			<div class="row">
				<form style="width: 500px; margin: auto" method="post">
					<h1 class="display">Anmälan</h1>
					För att kunna delta i tävlingen måste du först anmäla dig.
					<div class="form-group">
						<div class="submit-field">
							<input type="submit" value="Anmäl dig" class="raised">
						</div>
					</div>
				</form>
			</div>
		</article>
	</section>
{{ end }}
