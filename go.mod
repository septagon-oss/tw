module github.com/septagon-oss/tw

go 1.26

require github.com/septagon-dev/platformkit-design-system v0.0.0

require (
	github.com/septagon-oss/styleengine v0.0.0
	github.com/tdewolff/minify/v2 v2.24.13 // indirect
	github.com/tdewolff/parse/v2 v2.8.13 // indirect
)

replace github.com/septagon-dev/platformkit-design-system => ../../../frontend/platformkit-design-system

replace github.com/septagon-oss/styleengine => ../pk-styleengine
