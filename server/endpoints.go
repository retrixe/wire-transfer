package main

import (
	"encoding/json"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/retrixe/wire-transfer/core"
	"nhooyr.io/websocket"
)

func uploadEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{},
	})
	if err != nil {
		http.Error(w, errorJson("Failed to accept WebSocket connection!"), http.StatusBadRequest)
	}
	_, data, err := conn.Read(r.Context())
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "Failed to read from WebSocket!")
		return
	}
	var uploadRequest core.WsUploadRequest
	err = json.Unmarshal(data, &uploadRequest)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "Failed to parse JSON for upload request!")
		return
	}
	// TODO: Handle the rest.
	files.Store(gonanoid.MustGenerate(idChars, 8), File{})
}
