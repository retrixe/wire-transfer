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

## HTTP API

All request/response bodies use JSON.

### WS /upload

**Request:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file to upload. |
| `public_key` | `string` (optional) | The Ed25519 public key of the uploader. |
| `port` | `number` (optional) | The port on which the uploader is listening for direct transfers. |

**Response:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token which can be used to download the file. |

**Clientbound Request for Proxied Transfer:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `proxy_token` | `string` | The token to send to the server when connecting for proxied transfer. |

### GET /download/direct

**Request:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to download the file. |

**Response:**

This endpoint may return 404 if the file is not available, or 400 if the file is not available for direct transfer.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `ip` | `string` | The IP address of the uploader. |
| `port` | `number` | The port on which the uploader is listening for direct transfers. |

### GET /download/proxied

**Request:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to download the file. |

**Response:**

This endpoint may return 404 if the file is not available, or 400 if the file is not available for proxied transfer.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to send to the server when connecting for proxied transfer. |

## UDP Protocol API

The UDP protocol is used for transferring files.

### Packet Format

Packets are sent as a single UDP datagram. The first byte of the datagram is the packet type, followed by the packet data. When encryption is enabled, the contents of the packet are encrypted with AES-128-CBC.

### Type 0x00: Handshake

The handshake packet is sent by the downloader to the uploader or server to initiate a transfer. This is also sent by the uploader to the server when requested for a proxied transfer. **Note: Encryption fields are not present when this is sent by an uploader to a proxy server!**

**Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to download the file. In case of direct transfers, this is the download token. In case of proxied transfers, this is the proxying token, from which the download token can be inferred by the proxies/uploaders/downloaders. |
| `encryption` | `boolean` | Whether or not to enable encryption. This is not present when an uploader connects to a proxy server. |
| `shared_secret` | `string` (optional) | The shared secret to use for encryption, encrypted with the uploader's public key. Only present when `encryption` is set to `true`. |
| `iv` | `string` (optional) | The initialization vector to use for encryption, encrypted with the uploader's public key. Only present when `encryption` is set to `true`. |

### Type 0x01: File Data

WIP
