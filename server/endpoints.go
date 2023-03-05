package main

import (
	"encoding/json"
	"net/http"
	"time"

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
	id := gonanoid.MustGenerate(idChars, 8)
	file := File{
		Name:         uploadRequest.Name,
		Size:         uploadRequest.Size,
		Hash:         uploadRequest.Hash,
		CreationTime: time.Now(),
		PublicKey:    uploadRequest.PublicKey,
		Port:         uploadRequest.Port,
		Client:       conn,
	}
	files.Store(id, file)
	response, err := json.Marshal(core.WsUploadResponseSuccess{
		Success:  true,
		Token:    id,
		Precache: false,
	})
	if err != nil {
		conn.Close(websocket.StatusInternalError, "Failed to marshal JSON for upload response!")
		return
	}
	conn.Write(r.Context(), websocket.MessageText, response)

	// TODO: Keep monitoring the connection for disconnection. If the client disconnects, remove the client from the file.
	// TODO: Setup a timer to remove the file if the client doesn't connect within a certain amount of time.
}
