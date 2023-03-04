package main

import (
	"github.com/puzpuzpuz/xsync/v2"
)

type File struct{}

var files = xsync.NewMapOf[File]()
