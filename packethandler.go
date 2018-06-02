package main

import "fmt"
import "math/rand"

type PacketHandler struct {
	player Player
	packet *packet
}

func (handler *PacketHandler) sessionIDRequest() {
	uID := handler.packet.readByte()
	handler.player.session.uID = uID
	fmt.Printf("Read uID:%d\n", handler.player.session.uID)
	id := int64(rand.Uint32())<<32 + int64(rand.Uint32())
	handler.player.session.WriteLong(id)
}

func (handler *PacketHandler) HandlePacket() {
	switch handler.packet.opcode {
	case 32:
		handler.sessionIDRequest()
		break
	default:
		break
	}
}

func NewPacketHandler(player Player, packet *packet) PacketHandler {
	return PacketHandler{player, packet}
}
