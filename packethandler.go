package main

import "fmt"
import "math/rand"

type PacketHandler struct {
	player Player
}

func (handler *PacketHandler) loginRequest(packet *packet) {
	packet.readByte() // reconnecting := packet.readByte() == 1
	version := packet.readInt()
	fmt.Printf("Version:%d\n", version)
	packet.readInt()
	packet.readInt()
	packet.readInt()
	packet.readInt()
	userHash := packet.readLong()
	password := string(packet.data[packet.currentOffset:])
	handler.player.password = password
	fmt.Println(userHash)
	fmt.Println(password)
}

func (handler *PacketHandler) sessionIDRequest(packet *packet) {
	handler.player.session.uID = packet.readByte() & 0xFF
	id := int64(rand.Uint32())<<32 + int64(rand.Uint32())
	handler.player.session.id = id
	handler.player.session.WriteLong(id)
}

func (handler *PacketHandler) HandlePacket(packet *packet) {
	switch packet.opcode {
	case 32:
		handler.sessionIDRequest(packet)
		break
	case 0:
		handler.loginRequest(packet)
		break
	default:
		fmt.Printf("Unhandled packet[opcode:%d; length:%d]\n", packet.opcode, packet.length)
		break
	}
}

func NewPacketHandler(player Player) PacketHandler {
	return PacketHandler{player}
}
