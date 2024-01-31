package frontend

import (
	_ "embed"
)

//go:embed index.html
var HomeHTML string

//go:embed Gradient.js
var GradientJS []byte

//go:embed app.js
var AppJS []byte

//go:embed routes.js
var RoutesJS []byte

//go:embed style.css
var StyleCSS []byte

//go:embed logo.png
var LogoPNG []byte
