package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

// var hackBodyStyleRE = regexp.MustCompile(`^body \{[^}]*\}`)
// var hackBodyStyleRE = regexp.MustCompile(`body\s*{`)

var showCodeCache = make(map[string]string, 32)

// this is rather quick and dirty, should clean it up later, but it seems to work for now...
// Also, I think we're including all languages and styles even though we're only using a few,
// look at that also when cleaning this up.
func showCode(syntax string, code string) string {

	if ret, ok := showCodeCache[syntax+code]; ok {
		return ret
	}

	var buf bytes.Buffer
	err := quick.Highlight(&buf, code, syntax, "html", "monokailight")
	if err != nil {
		log.Print(err)
		return code
	}

	out := buf.String()
	// There's some crap in the output we need to hack out.
	// I'm sure it's possible to be smarter about this by digging into chroma more,
	// but doing this for now.
	out = strings.Replace(out, `<html>`, ``, 1)
	out = strings.Replace(out, `<body class="chroma">`, ``, 1)
	out = strings.Replace(out, `</body>`, ``, 1)
	out = strings.Replace(out, `</html>`, ``, 1)
	out = strings.Replace(out, `body {`, `nobody {`, 1)
	// out = hackBodyStyleRE.ReplaceAllLiteralString(out, ``)

	showCodeCache[syntax+code] = out

	return out
}
