module github.com/dreamph/cenery/test

go 1.24.0

replace (
	github.com/dreamph/cenery => ../
	github.com/dreamph/cenery/echo => ../engine/echo
	github.com/dreamph/cenery/fiber => ../engine/fiber
)

require (
	github.com/dreamph/cenery v1.0.1
	github.com/dreamph/cenery/echo v0.0.0
	github.com/dreamph/cenery/fiber v0.0.0
	github.com/goccy/go-json v0.10.5
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/labstack/echo/v4 v4.13.4
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.68.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.14.0 // indirect
)
