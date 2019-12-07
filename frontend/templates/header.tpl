{{ define "header" }}
<!DOCTYPE html>
<html lang="sv">
<head>
  <title>{{ template "title" . }}</title>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <script type="text/javascript" src="/static/kodsport.min.js"></script>
    <script type="text/javascript" src="/static/clipboard/clipboard.min.js"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/latest.js?config=TeX-MML-AM_HTMLorMML"></script>
  <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.blue-amber.min.css" />
  <link rel="stylesheet" href="/static/main.css">
  <link rel="stylesheet" href="/static/judge.css">
  <link rel="stylesheet" href="/static/editor.css">
  <link rel="stylesheet" href="/static/clipboard/primer-css/primer.css">
  <link rel="stylesheet" href="/static/clipboard/highlightjs/styles/github.css">
  <link rel="icon" href="/static/kodsport/favicon.png" type="image/png">
  <link href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:300,300i,400,400i,600,600i,700,700i" rel="stylesheet">
  <link href="https://fonts.googleapis.com/css?family=Roboto+Mono" rel="stylesheet">
  <link href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" rel="stylesheet">
  <link href="https://fonts.googleapis.com/css?family=Crimson+Text&display=swap" rel="stylesheet">
  <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
</head>
<body>
{{ end }}

{{ define "title" }}
  {{ if .C.Contest }}
    {{ .C.Contest.Title }} - Kodsport.Dev
  {{ else }}
    Kodsport.dev
  {{ end }}
{{ end }}
