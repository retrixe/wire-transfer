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

### Why not use UDP everywhere, instead of just the transfer protocol?

UDP is considerably harder to use, and HTTPS prevents any man-in-the-middle attacks by establishing a chain of trust. In future versions, we may reconsider this design decision, however, for now it works well and is easy to implement.

## HTTP API

All request/response bodies use JSON.

### WS /upload

**Request:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file to upload. |
| `size` | `number` | The size of the file, in bytes. |
| `hash` | `string` | The SHA256 hash of the file. |
| `public_key` | `string` (optional) | The Ed25519 public key of the uploader. |
| `port` | `number` (optional) | The port on which the uploader is listening for direct transfers. |
| `precache` | `boolean` (optional) | Hint to the server to try pre-cache the file. Pre-caching is not required for servers to implement. The reference implementation pre-caches by default. If you wish to avoid pre-caching as a client, set this to false explicitly. |

**Response:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `success` | `boolean` | Whether the request was successful. If the request failed, only the `error` field will be returned. |
| `error` | `string` (optional) | The error message, if the request was not successful. This is an implementation detail, and may not be always present. |
| `token` | `string` | The token which can be used to download the file. |
| `creation_time` | `number` | The time at which the file was created, in milliseconds since the UNIX epoch in UTC. |
| `precache` | `boolean` | Whether the file will be precached by the server. |

**Clientbound Request for Proxied Transfer:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `proxy_token` | `string` | The token to send to the server when connecting for proxied transfer. |

### GET /info

**Request Parameters:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to download the file. |

**Response Body:**

This endpoint may return 404 if the file is not available.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file. |
| `size` | `number` | The size of the file, in bytes. |
| `hash` | `string` | The SHA256 hash of the file. |
| `creation_time` | `number` | The time at which the file was created, in milliseconds since the UNIX epoch in UTC. |
| `available` | `boolean` | Whether the file is available. |
| `supports_direct` | `boolean` | Whether the file is available for direct transfer. |
| `supports_proxied` | `boolean` | Whether the file is available for proxied transfer. |

### GET /download/direct

**Request Parameters:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to download the file. |

**Response Body:**

This endpoint may return 404 if the file is not available, or 400 if the file is not available for direct transfer.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file. |
| `size` | `number` | The size of the file, in bytes. |
| `hash` | `string` | The SHA256 hash of the file. |
| `creation_time` | `number` | The time at which the file was created, in milliseconds since the UNIX epoch in UTC. |
| `ip` | `string` | The IP address of the uploader. |
| `port` | `number` | The port on which the uploader is listening for direct transfers. |

### GET /download/proxied

**Request Parameters:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | `string` | The token to send to the server when connecting for proxied transfer. |

**Response Body:**

This endpoint may return 404 if the file is not available, or 400 if the file is not available for proxied transfer.

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file. |
| `size` | `number` | The size of the file, in bytes. |
| `hash` | `string` | The SHA256 hash of the file. |
| `creation_time` | `number` | The time at which the file was created, in milliseconds since the UNIX epoch in UTC. |
| `token` | `string` | The token to download the file. |

## UDP Protocol API

The UDP protocol is used for transferring files.

### Packet Format

Packets are sent as a single UDP datagram. The first byte of the datagram is the packet type, followed by the packet data. When encryption is enabled, the contents of the packet are encrypted with AES-128-CBC.

### Data Types

| Type | Description |
| ---- | ----------- |
| `boolean` | A single byte, either `0x00` or `0x01`. |
| `string` | The length of the string encoded as a Protocol Buffer VarInt, followed by the string in UTF-8 encoded bytes. |
| `byte[]` | The length of the byte array encoded as a Protocol Buffer VarInt, followed by the byte array. If the length is specified, there is no VarInt prefixed. |

### Type 0x00: Handshake Request

The handshake packet is sent by the downloader to the uploader or server to initiate a transfer. This is also sent by the uploader to the server when requested for a proxied transfer. **Note: Encryption fields are not present when this is sent by an uploader to a proxy server!**

**Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `version` | `uint8` | The version of the protocol. This must be `1`. |
| `token` | `string` | The token to download the file. In case of direct transfers, this is the download token. In case of proxied transfers, this is the proxying token, from which the download token can be inferred by the proxies/uploaders/downloaders. |
| `encrypt` | `boolean` | Whether or not to enable encryption. This is not present when an uploader connects to a proxy server. |
| `shared_secret` | `byte[16]` (optional) | The shared secret to use for encryption, encrypted with the uploader's public key. Only present when `encrypt` is set to `true`. |
| `iv` | `byte[16]` (optional) | The initialization vector to use for encryption, encrypted with the uploader's public key. Only present when `encrypt` is set to `true`. |

### Type 0x00: Handshake Response

This packet is sent by the uploader in response to a handshake.

**Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `success` | `boolean` | Whether or not the handshake was successful. |
| `encrypt` | `boolean` | Whether or not to enable encryption. |

### Type 0x01: Request File Piece

This packet is sent by the downloader to the uploader to request a piece of the file. Each piece is 2048 KB in size.

**Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `piece` | `uint32` | The piece number to request. |
| `request_data` | `boolean` | Whether or not to request the data of the piece. |

### Type 0x01: File Piece Info

This packet is sent by the uploader to the downloader to send information about a piece of the file.

**Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `piece` | `uint32` | The piece number of the file. |
| `hash` | `byte[32]` | The SHA256 hash of the piece. |

### Type 0x02: File Piece Data

This packet is sent by the uploader to the downloader to send a piece of the file.

**Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `piece` | `uint32` | The piece number of the file. |
| `offset` | `uint32` | The offset of the data being received within the piece. |
| `data` | `byte[]` | The data of the piece. |

### 0x03: Close

This packet can be sent by the downloader to formally close the connection. This packet is not required to be sent, and the connection will be closed after a timeout (the reference implementation defaults to 30 seconds). However, using this saves resources.

**Packet Data:** N/A
