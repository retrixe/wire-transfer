package main

import (
	"time"

	"github.com/puzpuzpuz/xsync/v2"
)

type File struct {
	Name         string
	Size         int
	Hash         string
	CreationTime time.Time
	ExpiryTime   int
	// TODO: Any necessary fields for tracking uploaders and downloaders.
}

var files = xsync.NewMapOf[*File]()
