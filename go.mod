module github.com/teatak/pipe

go 1.13

require (
	github.com/teatak/cart v1.0.7
	github.com/teatak/config v0.0.0-00010101000000-000000000000
	github.com/teatak/riff v0.0.21 // indirect
)

replace (
	github.com/teatak/cart => ../cart
	github.com/teatak/config => ../config
)
