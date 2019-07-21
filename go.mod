module github.com/vugu/vugu-site

go 1.12

// replace github.com/vugu/vugu => ../vugu
replace github.com/vugu/vugu => ../vugu-stable-fixes

require (
	github.com/alecthomas/chroma v0.6.3
	github.com/vugu/vugu v0.0.0-00010101000000-000000000000
)
