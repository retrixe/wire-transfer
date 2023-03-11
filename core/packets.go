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
	ID   uint32
	Data []byte
}

func (packet *Packet) Serialize() []byte {
	data := make([]byte, 4+len(packet.Data))
	binary.LittleEndian.PutUint32(data[:4], packet.ID)
	copy(data[4:], packet.Data)
	return data
}

func ParsePacket(data []byte) (*Packet, error) {
	if len(data) < 4 {
		return nil, ErrInvalidPacket
	}
	packet := &Packet{}
	packet.ID = binary.LittleEndian.Uint32(data[:4])
	packet.Data = data[4:]
	return packet, nil
}

type InfoPacket struct {
	Version   uint8
	PublicKey []byte
}

func CreateInfoPacket(version uint8, publicKey []byte) *Packet {
	data := make([]byte, 2+len(publicKey))
	data[0] = version
	data[1] = boolToInt(publicKey != nil)
	if publicKey != nil {
		n := binary.PutVarint(data[2:], int64(len(publicKey)))
		copy(data[2+n:], publicKey)
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
	return &InfoPacket{
		Version:   packet.Data[0],
		PublicKey: packet.Data[1:],
	}, nil
}
