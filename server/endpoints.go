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

	// TODO: Validate all fields (limits, support for either/both direct/proxy transfer).

	id := gonanoid.MustGenerate(idChars, 8)
	file := &File{
		Name:         uploadRequest.Name,
		Size:         uploadRequest.Size,
		Hash:         uploadRequest.Hash,
		CreationTime: time.Now(),
		// TODO: Expiry time should be per-file with client hint (when protocol is added).
		ExpiryTime: config.DefaultFileExpiryTime,
		PublicKey:  uploadRequest.PublicKey,
		Port:       uploadRequest.Port,
		Client:     conn,
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

	// Keep monitoring the connection for disconnection.
	// If the client disconnects, remove the client from the file.
	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			file.Client = nil
			file.Reconnect = make(chan bool, 1)
			// Setup a timer to remove the file if the client doesn't reconnect within a certain amount of time.
			go func() {
				var response interface{}
				switch response {
				case file.ExpiryTime < 0:
					return
				case <-time.After(time.Duration(file.ExpiryTime) * time.Second):
					files.Delete(id)
					return
				case <-file.Reconnect:
					return
				}
			}()
			break
		}
	}
}
