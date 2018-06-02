package main

import "fmt"
import "math/rand"

type PacketHandler struct {
	player Player
}

func (handler *PacketHandler) sessionIDRequest(packet *packet) {
	handler.player.session.uID = packet.readByte() & 0xFF
	fmt.Printf("Read uID:%d\n", handler.player.session.uID)
	id := int64(rand.Uint32())<<32 + int64(rand.Uint32())
	handler.player.session.id = id
	handler.player.session.WriteLong(id)
}

func (handler *PacketHandler) HandlePacket(packet *packet) {
	switch packet.opcode {
	case 32:
		handler.sessionIDRequest(packet)
		break
	default:
		fmt.Printf("Unhandled packet[opcode:%d; length:%d]\n", packet.opcode, packet.length)
		break
	}
}

func NewPacketHandler(player Player) PacketHandler {
	return PacketHandler{player}
}
