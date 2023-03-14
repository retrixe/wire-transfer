# Wire Transfer Protocol Documentation

**Current version:** 0.x (in development)

This is a live document which serves as the reference for wire-transfer's protocol.

## Architecture

The Wire Transfer protocol is built around the concept of a central server and various clients which connect to it to upload and download files. The server is responsible for storing files, tracking their information, and issuing unique IDs to each file.

The protocol itself is only concerned with how clients can request information from servers (such as their public key, whether or not they are online, or information about any of the files they host) and how any 2 peers (server or client) can transfer a file between themselves efficiently, with the ability to locate and redownload corrupt parts of the tile. This flexibility allows for clients to reuse the same protocol to upload files to the server and for peer-to-peer transfers to work (the server can share info about where a client is listening for connections to download a certain file).

## To-dos

- Flesh out a transfer protocol which splits the file into pieces and allows for the client to request a piece (or all pieces) from the server. Hashes of each piece should also be part of the file metadata?

## Protocol

The UDP protocol is used for transferring files. Packets are sent as a single UDP datagram. The first byte of the datagram is the packet type, followed by the packet data. When encryption is enabled, the contents of the packet are encrypted with AES-128-CBC, for which the shared secret is derived from the ECDH (Curve25519) key pairs exchanged by the uploader and downloader.

Timeouts are left to implementations to decide. The reference implementation uses 30 seconds for downloads and 10 minutes for uploads to the server.

### Data Types

Optional data is prefixed with a boolean that is true when the data is present.

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
| `public_key` | `byte[]` (optional) | The server's ECDH (Curve25519) public key in PKIX, ASN.1 DER form. |
| `max_expiry_time` | `uint64` (optional) | The maximum time a file can be stored on the server, in milliseconds. |
| `max_file_size` | `uint64` (optional) | The maximum size a file can be, in bytes. |
| `info` | `string` (optional) | Extra information about the server. |

## Type 0x01: Handshake

This packet is used to establish a connection between two peers. The server will respond with a packet of type 0x01 as well.

**Request Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `version` | `uint8` | The version of the protocol. This must be `1`. |
| `public_key` | `byte[]` (optional) | The client's ECDH (Curve25519) public key in PKIX, ASN.1 DER form. |

**Response Packet Data:**

Identical to the [information packet's response](#type-0x00-information). After this, the client and server will use the shared secret derived from the ECDH key pairs to encrypt all packets sent to each other (if exchanged by both sides).

## Type 0x02: Close

This packet is used to close a connection between two peers. This packet may optionally contain an error message.

**Request Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `error` | `string` (optional) | An error message to send to the other peer. |

## Type 0x03: File Metadata

This packet is used to request information about a file from the server. The server will respond with a packet of type 0x01 as well.

**Request Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `id` | `string` | The ID of the file to request information about. |

**Response Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file. |
| `hash` | `byte[]` | The SHA-256 hash of the file. |
| `size` | `uint64` | The size of the file in bytes. |
| `available` | `boolean` | Whether or not the file is available for download. |
| `creation_time` | `uint64` | The time the file was created, in milliseconds since the Unix epoch. |
| `expiry_time` | `uint64` (optional) | The time the file will expire, in milliseconds since the Unix epoch. |
| `direct_download_ip` | `string` (optional) | The IP address of the server where the file can be downloaded directly. |
| `metadata` | `byte[]` (optional) | Extra metadata about the file. |

## Type 0x04: File Upload

This packet is used to upload a file to the server. The server will respond with a packet of type 0x02 as well.

**Request Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | `string` | The name of the file. |
| `hash` | `byte[]` | The SHA-256 hash of the file. |
| `size` | `uint64` | The size of the file in bytes. |
| `expiry_time` | `uint64` (optional) | The time the file will expire, in milliseconds since the Unix epoch. |
| `direct_download_ip` | `string` (optional) | The IP address of the server where the file can be downloaded directly. |
| `metadata` | `byte[]` (optional) | Extra metadata about the file. |

**Response Packet Data:**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `id` | `string` | The ID of the file. |
| `creation_time` | `uint64` | The time the file was created, in milliseconds since the Unix epoch. |
| `expiry_time` | `uint64` (optional) | The time the file will expire, in milliseconds since the Unix epoch. |
