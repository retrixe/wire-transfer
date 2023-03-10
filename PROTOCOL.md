# Wire Transfer Protocol Documentation

**Current version:** 0.x (in development)

This is a live document which serves as the reference for wire-transfer's protocol.

## Architecture

The Wire Transfer protocol is built around the concept of a central server and various clients which connect to it to upload and download files. The server is responsible for storing files, tracking their information, and issuing unique IDs to each file.

The protocol itself is only concerned with how clients can request information from servers (such as their public key, whether or not they are online, or information about any of the files they host) and how any 2 peers (server or client) can transfer a file between themselves efficiently, with the ability to locate and redownload corrupt parts of the tile. This flexibility allows for clients to reuse the same protocol to upload files to the server and for peer-to-peer transfers to work (the server can share info about where a client is listening for connections to download a certain file).

## To-dos

- Add a way for the client to upload a file to the server. File metadata includes name, hash, size, creation time, expiry time, IP and port where others could attempt the direct download, and extra metadata.
- Add a way to negotiate expiry time.
- Formally specify timeouts (30s, if file transfer is not complete in 30s, the client should request the file again).
- Flesh out a transfer protocol which splits the file into pieces and allows for the client to request a piece (or all pieces) from the server. Hashes of each piece should also be part of the file metadata?

## Protocol

The UDP protocol is used for transferring files. Packets are sent as a single UDP datagram. The first byte of the datagram is the packet type, followed by the packet data. When encryption is enabled, the contents of the packet are encrypted with AES-128-CBC, for which the shared secret is derived from the ECDH (Curve25519) key pairs exchanged by the uploader and downloader.

Timeouts are left to implementations to decide. The reference implementation uses 30 seconds for downloads and 10 minutes for uploads to the server.

### Data Types

| Type | Description |
| ---- | ----------- |
| `boolean` | A single byte, either `0x00` or `0x01`. |
| `string` | The length of the string encoded as a Protocol Buffer VarInt, followed by the string in UTF-8 encoded bytes. |
| `byte[]` | The length of the byte array encoded as a Protocol Buffer VarInt, followed by the byte array. If the length is specified, there is no VarInt prefixed. |

### Type 0x00: Information

This is a special packet used to request information from the server about itself. The server will respond with a packet of type 0x00 as well.

**Response Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `version` | `uint8` | The version of the protocol. This must be `1`. |
| `public_key` | `byte[]` (optional) | The server's Ed25519 public key. |
