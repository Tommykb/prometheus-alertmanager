package ui

const indexTmpl = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
        <title>Alertmanager</title>
        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css" integrity="sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ" crossorigin="anonymous">
        <script src="https://use.fontawesome.com/b7508bb100.js"></script>
    </head>
    <body>
        <!-- Your source after making -->
        <script src="{{ .ExternalURL }}/script.js"></script>
        <script>Elm.Main.embed(document.body, { production: true })</script>
    </body>
</html>`
