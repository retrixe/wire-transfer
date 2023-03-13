package core

import (
	"encoding/binary"
	"errors"
)

const ProtocolVersion = 0x01

var (
	ErrInvalidPacket = errors.New("invalid packet")
)

const (
	InfoPacketId = 0x00
)

type Packet struct {
	ID   uint8
	Data []byte
}

func (packet *Packet) Serialize() []byte {
	data := make([]byte, 1+len(packet.Data))
	data[0] = packet.ID
	copy(data[1:], packet.Data)
	return data
}

func ParsePacket(data []byte) (*Packet, error) {
	if len(data) < 1 {
		return nil, ErrInvalidPacket
	}
	packet := &Packet{}
	packet.ID = data[0]
	if len(data) >= 2 {
		packet.Data = data[1:]
	}
	return packet, nil
}

type InfoPacket struct {
	Version       uint8
	PublicKey     []byte
	MaxFileSize   uint64
	MaxExpiryTime uint64
	Info          string
}

func CreateInfoPacket(version uint8, publicKey []byte, maxFileSize *int, maxExpiryTime *int, info string) *Packet {
	data := make([]byte, 2)
	data[0] = version
	data[1] = boolToInt(publicKey != nil)
	if publicKey != nil {
		data = binary.AppendVarint(data, int64(len(publicKey)))
		data = append(data, publicKey...)
	}
	data = append(data, boolToInt(maxFileSize != nil))
	if maxFileSize != nil {
		data = binary.BigEndian.AppendUint64(data, uint64(*maxFileSize))
	}
	data = append(data, boolToInt(maxExpiryTime != nil))
	if maxExpiryTime != nil {
		data = binary.BigEndian.AppendUint64(data, uint64(*maxExpiryTime))
	}
	data = append(data, boolToInt(info != ""))
	if info != "" {
		data = binary.AppendVarint(data, int64(len([]byte(info))))
		data = append(data, []byte(info)...)
	}
	return &Packet{
		ID:   0x00,
		Data: data,
	}
}

func ParseInfoPacket(packet *Packet) (*InfoPacket, error) {
	if len(packet.Data) < 2 {
		return nil, ErrInvalidPacket
	}
	data := packet.Data

	version := data[0]
	data = data[1:]

	var publicKey []byte = nil
	if data[0] == 0x01 {
		if len(data) < 2 {
			return nil, ErrInvalidPacket
		}
		n, nLen := binary.Varint(data[1:])
		if nLen <= 0 || len(data) < 1+nLen+int(n) {
			return nil, ErrInvalidPacket
		}
		publicKey = data[1+nLen : 1+nLen+int(n)]
		data = data[1+nLen+int(n):]
	} else {
		data = data[1:]
	}

	if len(data) < 1 {
		return nil, ErrInvalidPacket
	}
	var maxFileSize uint64
	if data[0] == 0x01 {
		if len(data) < 9 {
			return nil, ErrInvalidPacket
		}
		maxFileSize = binary.BigEndian.Uint64(data[1:9])
		data = data[9:]
	} else {
		data = data[1:]
	}

	if len(data) < 1 {
		return nil, ErrInvalidPacket
	}
	var maxExpiryTime uint64
	if data[0] == 0x01 {
		if len(data) < 9 {
			return nil, ErrInvalidPacket
		}
		maxExpiryTime = binary.BigEndian.Uint64(data[1:9])
		data = data[9:]
	} else {
		data = data[1:]
	}

	if len(data) < 1 {
		return nil, ErrInvalidPacket
	}
	var info string = ""
	if data[0] == 0x01 {
		if len(data) < 2 {
			return nil, ErrInvalidPacket
		}
		n, nLen := binary.Varint(data[1:])
		if nLen <= 0 || len(data) < 1+nLen+int(n) {
			return nil, ErrInvalidPacket
		}
		info = string(data[1+nLen : 1+nLen+int(n)])
		// data = data[1+nLen+int(n):] } else { data = data[1:]
	}
	return &InfoPacket{
		Version:       version,
		PublicKey:     publicKey,
		MaxFileSize:   maxFileSize,
		MaxExpiryTime: maxExpiryTime,
		Info:          info,
	}, nil
}
