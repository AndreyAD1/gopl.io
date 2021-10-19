package main

import "html/template"

var MainPage = `<!DOCTYPE html>
<html>
<head>
<title>Web Calculator</title>
</head>
<body>

<h1>Web Calculator</h1>
<h2>Enter an Expression</h2>
<form action="/calculate">
<label for="expression">Math expression:</label><br>
<input type="text" id="expression" name="expression"><br>
<input type="submit" value="Submit">
</form> 

</body>
</html>`

var ResultPageTemplate = template.Must(template.New("ResultPage").Parse(`
<!DOCTYPE html>
<html>
<head>
<title>Web Calculator</title>
</head>
<body>

<h1>Web Calculator</h1>
<h2>Expression Result</h2>
<p>Expression: {{.Expression}}</p>
<p>Result: {{.Result}}</p>

</body>
</html>`))

var ErrorPage = `
<!DOCTYPE html>
<html>
<head>
<title>Web Calculator/title>
</head>
<body>

<h1>Internal Server Error</h1>
<p>Something went wrong.</p>

</body>
</html>
`