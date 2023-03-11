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
	Version   uint8
	PublicKey []byte
}

func CreateInfoPacket(version uint8, publicKey []byte) *Packet {
	data := make([]byte, 2)
	data[0] = version
	data[1] = boolToInt(publicKey != nil)
	if publicKey != nil {
		data = binary.AppendVarint(data, int64(len(publicKey)))
		data = append(data, publicKey...)
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
	var publicKey []byte = nil
	if packet.Data[1] == 0x01 {
		if len(packet.Data) < 3 {
			return nil, ErrInvalidPacket
		}
		n, nLen := binary.Varint(packet.Data[2:])
		if nLen <= 0 || len(packet.Data) < 2+nLen+int(n) {
			return nil, ErrInvalidPacket
		}
		publicKey = packet.Data[2+nLen : 2+nLen+int(n)]
	}
	return &InfoPacket{
		Version:   packet.Data[0],
		PublicKey: publicKey,
	}, nil
}
