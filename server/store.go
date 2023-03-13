package main

import (
	"net"
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

type Connection struct {
	Addr     *net.UDPAddr
	LastSeen time.Time
	// Send true when terminating the connection, false when simply indicating last seen.
	TimeoutChan  chan bool
	SharedSecret []byte
}

var connections = xsync.NewMapOf[*Connection]()
