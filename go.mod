module github.com/omakoto/gocmds

go 1.14

//replace github.com/omakoto/go-common => ../go-common
//replace github.com/omakoto/bashcomp => ../bashcomp

require (
	github.com/omakoto/bashcomp v0.0.0-20160616051942-c0902e9f1d64
	github.com/omakoto/go-common v0.0.0-20200711204306-19446ee8d4ef
	github.com/pborman/getopt v0.0.0-20190409184431-ee0cd42419d3
	github.com/pkg/errors v0.8.1
	github.com/shopspring/decimal v0.0.0-20190905144223-a36b5d85f337
	github.com/stretchr/testify v1.4.0
	github.com/ungerik/go-dry v0.0.0-20180411133923-654ae31114c8
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/sys v0.0.0-20190927073244-c990c680b611 // indirect
)
