package main

import "math/rand"

type PacketHandler struct {
	client *Session
	packet *packet
}

func (handler *PacketHandler) sessionIDRequest() {
	id := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	handler.client.conn.Write([]byte{byte(id >> 56), byte(id >> 48), byte(id >> 40), byte(id >> 32), byte(id >> 24), byte(id >> 16), byte(id >> 8), byte(id & 0xFF)})
}

func (handler *PacketHandler) HandlePacket() {
	switch handler.packet.Opcode {
	case 32:
		handler.sessionIDRequest()
		break
	default:
		break
	}
}

func NewPacketHandler(client *Session, packet *packet) PacketHandler {
	ph := PacketHandler{client, packet}
	return ph
}
