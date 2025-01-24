module github.com/dreamph/cenery/fiber

go 1.23

replace (
	github.com/dreamph/cenery => ../
	github.com/dreamph/cenery/engine/echo => ../engine/echo
	github.com/dreamph/cenery/engine/fiber => ../engine/fiber
)

require (
	github.com/dreamph/cenery v1.0.0
	github.com/dreamph/cenery/engine/echo v0.0.0-00010101000000-000000000000
	github.com/dreamph/cenery/engine/fiber v0.0.0-00010101000000-000000000000
	github.com/goccy/go-json v0.10.4
	github.com/gofiber/fiber/v2 v2.52.6
	github.com/labstack/echo/v4 v4.13.3
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.9.0 // indirect
)
