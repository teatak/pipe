module github.com/teatak/pipe

go 1.13

require (
	github.com/teatak/config v0.0.0-00010101000000-000000000000
	github.com/teatak/riff v0.0.21
)

replace (
	github.com/teatak/cart => ../cart
	github.com/teatak/config => ../config
)
