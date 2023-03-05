package main

import (
	"time"

	"github.com/puzpuzpuz/xsync/v2"
	"nhooyr.io/websocket"
)

type File struct {
	Name         string
	Size         int
	Hash         string
	CreationTime time.Time
	ExpiryTime   int
	// Only present if encryption is supported by the client.
	PublicKey string
	// Only present if direct transfers are supported by the client.
	Port int
	// Represents the connected client, if the client is connected.
	Client *websocket.Conn
	// This channel is set when the client is disconnected, and indicates any reconnection.
	Reconnect chan bool
	// TODO: Any necessary fields for tracking downloaders.
}

var files = xsync.NewMapOf[*File]()
