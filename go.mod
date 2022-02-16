module github.com/maldan/gam

go 1.18

// replace github.com/maldan/go-restserver => ../../go_lib/restserver
replace github.com/maldan/go-cmhp => ../../go_lib/cmhp

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.12.0
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/maldan/go-cmhp v0.0.19
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/schollz/progressbar/v3 v3.8.1
	golang.org/x/mod v0.4.2
)

require (
	github.com/aws/aws-sdk-go v1.40.6 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210603125802-9665404d3644 // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
)
