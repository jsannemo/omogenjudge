{{ define  "problem_sample" }}
  <div class="samplebox">
    <div class="samplebox-header">
      Exempelindata
      <span class="copy-sample-btn tooltipable" title="Copy to Clipboard">
        <i class="material-icons">file_copy</i>
      </span>
    </div>
    <pre class="samplebox-data">{{ .InputFile.FileString }}</pre>
    <div class="samplebox-header">
      Exempelutdata
      <span class="copy-sample-btn tooltipable" title="Copy to Clipboard">
        <i class="material-icons">file_copy</i>
      </span>
    </div>
    <pre class="samplebox-data">{{ .OutputFile.FileString }}</pre>
  </div>
{{ end }}
