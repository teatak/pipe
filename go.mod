module github.com/teatak/pipe

go 1.13

require (
	github.com/elazarl/go-bindata-assetfs v1.0.1 // indirect
	github.com/teatak/cart v0.0.0-00010101000000-000000000000
	github.com/teatak/config v0.0.0-00010101000000-000000000000
)

replace (
	github.com/teatak/cart => ../cart
	github.com/teatak/config => ../config
)
