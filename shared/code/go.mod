module code

go 1.25.0

replace common => ../common

replace auth => ../../auth-service

replace pb => ../protobuf

require (
	github.com/gdamore/tcell/v2 v2.8.1
	github.com/jinzhu/copier v0.4.0
	github.com/rivo/tview v0.0.0-20250625164341-a4a78f1e05cb
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
