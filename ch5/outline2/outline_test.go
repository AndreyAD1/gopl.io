package main

import (
	"testing"
	"strings"
)

func TestOutlineHTML(t *testing.T) {
	var testHTML = []struct {
		rawHtml string
		expectedPrint string
	}{
		{
			"<!DOCTYPE html><html><body><h1>My First Heading</h1><p>My first paragraph.</p></body></html>",
			`
<!DOCTYPE html>
	<html>
		<body>	
			<h1>
				My First Heading
			</h1>	
			<p>
				My first paragraph.
			</p>
		</body>
	</html>
`,
		},
	}
	for _, test := range testHTML {
		bodyReader := strings.NewReader(test.rawHtml)
		err := outline(bodyReader)
		if err != nil {
			t.Errorf("outline(%q) causes the error %s", test.rawHtml, err)
		}
	}
}