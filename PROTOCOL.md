# Wire Transfer Protocol Documentation

**Current version:** 0.x (in development)

This is a live document which serves as the reference for wire-transfer's protocol.

## Architecture

Wire Transfer relies upon a central server and various clients which connect to it. The server is responsible for maintaining a list of clients and their current status, and for coordinating the transfer of files between clients. The server may optionally implement caching on its end, to reduce the amount of data that needs to be transferred (however, the reference implementation does not currently implement this, although it may in the future).

To transfer a file, a client first connects to the server and requests to upload a file via WebSocket at `/upload`. This request contains the file name, an Ed25519 public key (optionally, to support encryption), and the port on which the client is listening for requesting clients (if wishing to support direct transfer). The server responds to this request with a token which can be shared with other clients to download the file.

While the uploader is connected via WebSocket, a requester can request the file using two modes, direct transfer and proxied transfer. The reference implementation prefers direct transfer, falling back to proxied transfer when direct transfer is unavailable, but clients can prefer/support either as they see fit.

### Direct Transfer

In direct transfer mode, the requesting client makes a request at `/download/direct`. If direct transfer is supported by the uploader, the server responds with the uploader's IP address, port and public key. The requester then connects to the uploader directly via UDP, and the uploader sends the file to the requester using the Wire Transfer UDP protocol. The uploader and downloader are responsible for encrypting and decrypting the file.

### Proxied Transfer

This mode is primarily aimed at users which don't support or enable port forwarding.

In proxied transfer mode, the requesting client makes a request at `/download/proxied`. If proxied transfer is supported by the uploader, the server generates a token and sends it to the uploader. The uploader then connects via UDP and sends the token, following which the server responds to the requester's HTTP request with the temporary token. The requester then connects to the server via UDP and sends the token. The server then relays the file between the uploader and the requester. The uploader and downloader are responsible for encrypting and decrypting the file, and the server is responsible for relaying the file.

## UDP Transfer Protocol

WIP
