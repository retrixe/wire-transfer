package core

import (
	"encoding/binary"
	"errors"
)

const ProtocolVersion = 0x01

var (
	ErrInvalidPacket = errors.New("invalid packet")
)

type Packet struct {
	ID     uint32
	Data   []byte
	Parsed interface{}
}

type InformationPacket struct {
	Version   uint8
	PublicKey []byte
}

/* type HandshakeRequest struct {
	Version      uint8
	Token        string
	Encrypt      bool
	SharedSecret []byte
	Iv           []byte
}

type HandshakeResponse struct {
	Success bool
	Encrypt bool
}

type RequestFilePiece struct {
	Piece       uint32
	RequestData bool
}

type FilePieceInfo struct {
	Piece uint32
	Hash  []byte
}

type FilePieceData struct {
	Piece  uint32
	Offset uint32
	Data   []byte
} */

func (packet *Packet) Serialize() []byte {
	data := make([]byte, 4+len(packet.Data))
	binary.LittleEndian.PutUint32(data[:4], packet.ID)
	copy(data[4:], packet.Data)
	return data
}

func CreateInformationPacket(version uint8, publicKey []byte) *Packet {
	data := make([]byte, 2+len(publicKey))
	data[0] = version
	if publicKey != nil {
		data[1] = 0x01
		copy(data[1:], publicKey)
	} else {
		data[1] = 0x00
	}
	return &Packet{
		ID:   0x00,
		Data: data,
	}
}

/* func CreateHandshakeRequest(token string, encrypt bool, sharedSecret []byte, iv []byte) *Packet {
	data := make([]byte, 3+len(token))
	data[0] = ProtocolVersion
	n := binary.PutVarint(data[1:], int64(len(token)))
	copy(data[1+n:], []byte(token))
	data[1+n+len(token)] = boolToInt(encrypt)
	if encrypt {
		copy(data[1+1+n+len(token):], sharedSecret)
		copy(data[16+1+1+n+len(token):], iv)
	}
	return &Packet{
		ID:   0x00,
		Data: data,
	}
}

func CreateHandshakeResponse(success bool, encrypt bool) *Packet {
	data := make([]byte, 2)
	data[0] = boolToInt(success)
	data[1] = boolToInt(encrypt)
	return &Packet{
		ID:   0x00,
		Data: data,
	}
} */

func ParsePacket(data []byte) (*Packet, error) {
	if len(data) < 4 {
		return nil, errors.New("invalid packet")
	}
	packet := &Packet{}
	packet.ID = binary.LittleEndian.Uint32(data[:4])
	packet.Data = data[4:]
	return packet, nil
}

/* func ParsePacketData(packet *Packet, from PacketFrom) error {
	if from == PacketFromDownloader {
		switch packet.ID {
		case 0x00:
			return parseHandshakeRequest(packet)
		case 0x01:
			return parseRequestFilePiece(packet)
		}
	} else if from == PacketFromUploader {
		switch packet.ID {
		case 0x00:
			return parseHandshakeResponse(packet)
		case 0x01:
			return parseFilePieceInfo(packet)
		case 0x02:
			return parseFilePieceData(packet)
		}
	}
	return ErrInvalidPacket
}

func parseHandshakeRequest(packet *Packet) error {
	if len(packet.Data) < 3 { // Allows for 1 byte version, 1 byte token length, 1 byte encrypt
		return ErrInvalidPacket
	}
	packet.Parsed = &HandshakeRequest{}
	index := 0
	packet.Parsed.(*HandshakeRequest).Version = packet.Data[index]
	index += 1

	tokenLength, n := binary.Varint(packet.Data[index:])
	index += n

	if len(packet.Data) < index+int(tokenLength) {
		return ErrInvalidPacket
	}
	packet.Parsed.(*HandshakeRequest).Token = string(packet.Data[index : index+int(tokenLength)])
	index += int(tokenLength)

	packet.Parsed.(*HandshakeRequest).Encrypt = packet.Data[index] == 0x01
	index += 1

	if packet.Parsed.(*HandshakeRequest).Encrypt {
		if len(packet.Data) < index+16 {
			return ErrInvalidPacket
		}
		packet.Parsed.(*HandshakeRequest).SharedSecret = packet.Data[index : index+16]
		index += 16

		if len(packet.Data) < index+16 {
			return ErrInvalidPacket
		}
		packet.Parsed.(*HandshakeRequest).Iv = packet.Data[index : index+16]
		index += 16
	}
	return nil
}

func parseHandshakeResponse(packet *Packet) error {
	if len(packet.Data) < 2 {
		return ErrInvalidPacket
	}
	packet.Parsed = &HandshakeResponse{}
	packet.Parsed.(*HandshakeResponse).Success = packet.Data[0] == 0x01
	packet.Parsed.(*HandshakeResponse).Encrypt = packet.Data[1] == 0x01
	return nil
}

func parseRequestFilePiece(packet *Packet) error {
	if len(packet.Data) < 5 {
		return ErrInvalidPacket
	}
	packet.Parsed = &RequestFilePiece{}
	packet.Parsed.(*RequestFilePiece).Piece = binary.LittleEndian.Uint32(packet.Data[:4])
	packet.Parsed.(*RequestFilePiece).RequestData = packet.Data[4] == 0x01
	return nil
}

func parseFilePieceInfo(packet *Packet) error {
	if len(packet.Data) < 36 {
		return ErrInvalidPacket
	}
	packet.Parsed = &FilePieceInfo{}
	packet.Parsed.(*FilePieceInfo).Piece = binary.LittleEndian.Uint32(packet.Data[:4])
	packet.Parsed.(*FilePieceInfo).Hash = packet.Data[4 : 4+32]
	return nil
}

func parseFilePieceData(packet *Packet) error {
	if len(packet.Data) < 9 {
		return ErrInvalidPacket
	}
	packet.Parsed = &FilePieceData{}
	packet.Parsed.(*FilePieceData).Piece = binary.LittleEndian.Uint32(packet.Data[:4])
	packet.Parsed.(*FilePieceData).Offset = binary.LittleEndian.Uint32(packet.Data[4:8])
	dataLength, n := binary.Varint(packet.Data[8:])
	if int(dataLength) != len(packet.Data[8+n:]) {
		return ErrInvalidPacket
	}
	packet.Parsed.(*FilePieceData).Data = packet.Data[8+n:]
	return nil
}

func boolToInt(b bool) uint8 {
	if b {
		return 0x01
	}
	return 0x00
} */
