{{ define  "problem_sample" }}
  <div class="samplebox">
    <div class="samplebox-header">Exempelindata</div>
    <pre class="samplebox-data">{{ .InputFile.FileString }}</pre>
    <div class="samplebox-header">Exempelutdata</div>
    <pre class="samplebox-data">{{ .OutputFile.FileString }}</pre>
  </div>
{{ end }}
